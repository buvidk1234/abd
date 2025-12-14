package service

import (
	"backend/internal/model"
	"backend/internal/pkg/cache/cachekey"
	"backend/internal/pkg/cache/redis"
	"backend/internal/pkg/constant"
	"backend/internal/pkg/snowflake"
	"context"
	"errors"
	"log"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageService struct {
	db           *gorm.DB
	seqConvCache *redis.SeqConversationCacheRedis
	seqUserCache *redis.SeqUserCacheRedis
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{
		db:           db,
		seqConvCache: redis.NewSeqConversationCacheRedis(db, redis.GetRDB()),
		seqUserCache: redis.NewSeqUserCacheRedis(db, redis.GetRDB()),
	}
}

type SendMessageReq struct {
	SenderID int64  `json:"sender_id,string"`
	ConvType int32  `json:"conv_type" binding:"required"`
	TargetID int64  `json:"target_id,string" binding:"required"`
	MsgType  int32  `json:"msg_type" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// Deprecated: use im/distributor instead
func (s *MessageService) SendMessage(ctx context.Context, req SendMessageReq) error {
	conversationID := GetConversationID(req.ConvType, req.SenderID, req.TargetID)

	// TODO: 创建好友时或加入群聊时调用，初始化会话记录
	// s.InitConversationForUser(ctx, req)

	var newSeq int64
	var msg model.Message
	// perform read->increment->update and message create in a single transaction
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var seqConversation model.SeqConversation
		// lock the seq row to avoid concurrent increments
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", conversationID).FirstOrCreate(&seqConversation).Error; err != nil {
			return err
		}

		newSeq = seqConversation.MaxSeq + 1

		if err := tx.Model(&model.SeqConversation{}).Where("id = ?", conversationID).Update("max_seq", newSeq).Error; err != nil {
			return err
		}

		msg = model.Message{
			ID:             snowflake.GenID(),
			ConversationID: conversationID,
			Seq:            newSeq,
			SenderID:       req.SenderID,
			MsgType:        req.MsgType,
			Content:        req.Content,
		}

		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// async update user timeline
	go func() {
		switch req.ConvType {
		case constant.SingleChatType:
			err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				var seqUser model.SeqUser

				// FOR UPDATE + FirstOrCreate 合并写法
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
					FirstOrCreate(&seqUser, "id = ?", req.TargetID).Error; err != nil {
					return err
				}

				newSeqUser := seqUser.MaxSeq + 1

				if err := tx.Model(&model.SeqUser{}).Where("id = ?", req.TargetID).Update("max_seq", newSeqUser).Error; err != nil {
					return err
				}

				userTimeline := model.UserTimeline{
					OwnerID:        req.TargetID,
					Seq:            newSeqUser,
					ConversationID: conversationID,
					MsgID:          msg.ID,
					RefMsgSeq:      newSeq,
					MsgType:        req.MsgType,
					SenderID:       req.SenderID,
					Snapshot:       req.Content, // TODO：生成消息摘要
				}

				if err := tx.Create(&userTimeline).Error; err != nil {
					return err
				}

				return nil
			})

			// goroutine 中错误无法返回，需要 channel 或 log
			if err != nil {
				log.Printf("[timeline error] %v", err)
			}
		case constant.GroupChatType:
			/*
				用 Redis 批量生成 1000 个 Seq。
				在内存里构建 1000 个 UserTimeline 对象。
				只开 1 个事务。
				执行一次 tx.CreateInBatches(timelines, 100)
			*/
			log.Printf("[timeline] group message timeline update not implemented")
		default:
			log.Printf("err conv type")
			return
		}
	}()

	return nil
}

type PullSpecifiedConvReq struct {
	UserID  int64  `json:"user_id,string"`
	ConvID  string `json:"conv_id"`
	ConvSeq int64  `json:"conv_seq"`
}
type PullSpecifiedConvResp struct {
	Messages []model.Message `json:"messages"`
}

func (s *MessageService) PullSpecifiedConv(ctx context.Context, req PullSpecifiedConvReq) (PullSpecifiedConvResp, error) {

	var msgs []model.Message
	if err := s.db.WithContext(ctx).Find(&msgs, "conversation_id = ? AND seq >= ? ORDER BY seq ASC", req.ConvID, req.ConvSeq).Error; err != nil {
		return PullSpecifiedConvResp{}, err
	}
	return PullSpecifiedConvResp{Messages: msgs}, nil
}

type PullConvListReq struct {
	UserID  int64 `json:"user_id,string"`
	UserSeq int64 `json:"user_seq"`
}
type PullConvListResp struct {
	PullMsgs map[string][]model.Message `json:"pull_msgs"`
}

func (s *MessageService) PullConvList(ctx context.Context, req PullConvListReq) (PullConvListResp, error) {
	var userTimelines []model.UserTimeline
	if err := s.db.WithContext(ctx).Find(&userTimelines, "owner_id = ? AND seq >= ? ORDER BY seq ASC", req.UserID, req.UserSeq).Error; err != nil {
		return PullConvListResp{}, err
	}
	pullMsgs := make(map[string][]model.Message)
	for _, timeline := range userTimelines {
		msgAbstract := model.Message{
			ConversationID: timeline.ConversationID,
			ID:             timeline.MsgID,
			Seq:            timeline.RefMsgSeq,
			MsgType:        timeline.MsgType,
			SenderID:       timeline.SenderID,
			Content:        timeline.Snapshot,
		}
		pullMsgs[timeline.ConversationID] = append(pullMsgs[timeline.ConversationID], msgAbstract)
	}
	return PullConvListResp{PullMsgs: pullMsgs}, nil
}

func (s *MessageService) DeleteConversation(ctx context.Context, userID int64, conversationID string) error {
	var seqConversation model.SeqConversation

	// 查询 SeqConversation
	if err := s.db.WithContext(ctx).Find(&seqConversation, "conversation_id = ?", conversationID).Error; err != nil {
		return err
	}

	// 更新 Conversation 的 read_seq 和 max_seq
	if err := s.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Updates(map[string]interface{}{
			"read_seq": seqConversation.MaxSeq,
			"max_seq":  seqConversation.MaxSeq,
		}).Error; err != nil {
		return err
	}

	return nil
}

type GetMaxSeqReq struct {
	UserID int64 `json:"user_id,string"`
}
type GetMaxSeqResp struct {
	MaxSeqs map[string]int64 `json:"max_seqs"`
	MinSeqs map[string]int64 `json:"min_seqs"`
}

func (s *MessageService) GetMaxSeq(ctx context.Context, req GetMaxSeqReq) (GetMaxSeqResp, error) {

	conversationIDs, err := redis.GetCache(cachekey.GetConversationIDsKey(strconv.FormatInt(req.UserID, 10)), func() ([]string, error) {
		var ids []string
		if err := s.db.WithContext(ctx).Model(&model.Conversation{}).Where("owner_id = ?", req.UserID).Pluck("conversation_id", &ids).Error; err != nil {
			return nil, err
		}
		return ids, nil
	}, redis.ExpireTime)

	if err != nil {
		return GetMaxSeqResp{}, err
	}

	maxSeqs, err := s.seqConvCache.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return GetMaxSeqResp{}, err
	}
	// avoid pulling messages from sessions with a large number of max seq values of 0
	for conversationID, seq := range maxSeqs {
		if seq == 0 {
			delete(maxSeqs, conversationID)
		}
	}
	return GetMaxSeqResp{MaxSeqs: maxSeqs}, nil
}

type SeqRange struct {
	ConversationID string `json:"conversation_id"`
	Begin          int64  `json:"begin"`
	End            int64  `json:"end"`
	Num            int64  `json:"num"`
}

type PullOrder int

const (
	PullOrderAsc  PullOrder = 1
	PullOrderDesc PullOrder = 2
)

type PullMsgs struct {
	Msgs   []*model.Message `json:"msgs"`
	IsEnd  bool             `json:"is_end"`
	EndSeq int64            `json:"end_seq"`
}

type PullMessageBySeqsReq struct {
	UserID    int64       `json:"user_id,string"`
	SeqRanges []*SeqRange `json:"seq_ranges"`
	Order     PullOrder   `json:"order"`
}

type PullMessageBySeqsResp struct {
	Msgs             map[string]*PullMsgs `json:"msgs"`
	NotificationMsgs map[string]*PullMsgs `json:"notification_msgs"`
}

func (s *MessageService) PullMessageBySeqs(ctx context.Context, req PullMessageBySeqsReq) (PullMessageBySeqsResp, error) {
	resp := PullMessageBySeqsResp{
		Msgs:             make(map[string]*PullMsgs),
		NotificationMsgs: make(map[string]*PullMsgs),
	}

	for _, seqRange := range req.SeqRanges {
		log.Printf("PullMessageBySeqs processing conversationID: %v, begin: %v, end: %v, num: %v", seqRange.ConversationID, seqRange.Begin, seqRange.End, seqRange.Num)
		conversation, err := redis.GetCache(cachekey.GetConversationKey(strconv.FormatInt(req.UserID, 10), seqRange.ConversationID), func() (model.Conversation, error) {
			var conv model.Conversation
			if err := s.db.WithContext(ctx).Where("owner_id = ? AND conversation_id = ?", req.UserID, seqRange.ConversationID).First(&conv).Error; err != nil {
				return model.Conversation{}, err
			}
			return conv, nil
		}, redis.ExpireTime)
		if err != nil {
			log.Printf("PullMessageBySeqs get conversation error: %v, conversationID: %v", err, seqRange.ConversationID)
			continue
		}
		minSeq, maxSeq, msgs, err := s.getMsgBySeqsRange(ctx, req.UserID, seqRange.ConversationID, seqRange.Begin, seqRange.End, seqRange.Num, conversation.MaxSeq)
		if err != nil {
			log.Printf("PullMessageBySeqs get messages error: %v, conversationID: %v", err, seqRange.ConversationID)
			continue
		}
		var isEnd bool
		switch req.Order {
		case PullOrderAsc:
			isEnd = (maxSeq <= seqRange.End)
		case PullOrderDesc:
			isEnd = (minSeq >= seqRange.Begin)
		}
		if len(msgs) == 0 {
			log.Printf("PullMessageBySeqs no messages found, conversationID: %v, begin: %v, end: %v", seqRange.ConversationID, seqRange.Begin, seqRange.End)
			continue
		}
		resp.Msgs[seqRange.ConversationID] = &PullMsgs{
			Msgs:  msgs,
			IsEnd: isEnd,
		}
	}

	return resp, nil
}

func (s *MessageService) getMsgBySeqsRange(ctx context.Context, userID int64, conversationID string, begin, end, num, userMaxSeq int64) (int64, int64, []*model.Message, error) {
	userMinSeq, err := s.seqUserCache.GetSeqUserMinSeq(ctx, userID, conversationID)
	if err != nil {
		return 0, 0, nil, err
	}
	minSeq, err := s.seqConvCache.GetMinSeq(ctx, conversationID)
	if err != nil {
		return 0, 0, nil, err
	}
	if userMinSeq > minSeq {
		minSeq = userMinSeq
	}
	// "minSeq" represents the startSeq value that the user can retrieve.
	if minSeq > end {
		log.Printf("getMsgBySeqsRange no messages to pull, userMinSeq: %v, conMinSeq: %v, begin: %v, end: %v", userMinSeq, minSeq, begin, end)
		return 0, 0, nil, nil
	}
	maxSeq, err := s.seqConvCache.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return 0, 0, nil, err
	}
	log.Printf("getMsgBySeqsRange before adjust: minSeq=%d, maxSeq=%d, userMinSeq=%d, userMaxSeq=%d, begin=%d, end=%d", minSeq, maxSeq, userMinSeq, userMaxSeq, begin, end)
	if userMaxSeq != 0 {
		if userMaxSeq < maxSeq {
			maxSeq = userMaxSeq
		}
	}
	// "maxSeq" represents the endSeq value that the user can retrieve.

	if begin < minSeq {
		begin = minSeq
	}
	if end > maxSeq {
		end = maxSeq
	}
	// "begin" and "end" represent the actual startSeq and endSeq values that the user can retrieve.
	if end < begin {
		return 0, 0, nil, errors.New("seq end < begin")
	}
	var seqs []int64
	if end-begin+1 <= num {
		for i := begin; i <= end; i++ {
			seqs = append(seqs, i)
		}
	} else {
		for i := end - num + 1; i <= end; i++ {
			seqs = append(seqs, i)
		}
	}
	successMsgs, err := s.GetMessageBySeqs(ctx, conversationID, userID, seqs)
	if err != nil {
		return 0, 0, nil, err
	}
	return minSeq, maxSeq, successMsgs, nil
}

func (s *MessageService) GetMessageBySeqs(ctx context.Context, conversationID string, userID int64, seqs []int64) ([]*model.Message, error) {
	var keys []string
	var keyseqMap = make(map[string]int64)
	for _, seq := range seqs {
		key := cachekey.GetMsgCacheKey(conversationID, seq)
		keyseqMap[key] = seq
		keys = append(keys, key)
	}
	msgs, err := redis.BatchGetCache(keys, func(missingKeys []string) (map[string]model.Message, error) {
		var messages []model.Message
		if err := s.db.WithContext(ctx).Find(&messages, "conversation_id = ? AND seq IN ?", conversationID, func() []int64 {
			var ms []int64
			for _, k := range missingKeys {
				ms = append(ms, keyseqMap[k])
			}
			return ms
		}()).Error; err != nil {
			return nil, err
		}
		result := make(map[string]model.Message)
		for _, msg := range messages {
			result[cachekey.GetMsgCacheKey(conversationID, msg.Seq)] = msg
		}
		return result, nil
	}, redis.ExpireTime)
	if err != nil {
		return nil, err
	}
	var result []*model.Message
	for _, key := range keys {
		if msg, exists := msgs[key]; exists {
			m := msg // create a new variable to take the address
			result = append(result, &m)
		}
	}
	return result, nil
}

// GetSeqMessage
type ConversationSeqs struct {
	ConversationID string  `json:"conversation_id"`
	Seqs           []int64 `json:"seqs"`
}

type GetSeqMessageReq struct {
	UserID        int64               `json:"user_id,string"`
	Conversations []*ConversationSeqs `json:"conversations"`
	Order         PullOrder           `json:"order"`
}

type GetSeqMessageResp struct {
	Msgs map[string]*PullMsgs `json:"msgs"`
}

func (s *MessageService) GetSeqMessage(ctx context.Context, req GetSeqMessageReq) (GetSeqMessageResp, error) {
	resp := GetSeqMessageResp{
		Msgs: make(map[string]*PullMsgs),
	}
	for _, conv := range req.Conversations {
		isEnd, endSeq, msgs, err := s.GetMessagesBySeqWithBounds(ctx, req.UserID, conv.ConversationID, conv.Seqs, req.Order)
		if err != nil {
			log.Printf("GetSeqMessage error: %v, conversationID: %v", err, conv.ConversationID)
			continue
		}
		resp.Msgs[conv.ConversationID] = &PullMsgs{
			Msgs:   msgs,
			IsEnd:  isEnd,
			EndSeq: endSeq,
		}
	}
	return resp, nil
}
func (s *MessageService) GetMessagesBySeqWithBounds(ctx context.Context, userID int64, conversationID string, seqs []int64, pullOrder PullOrder) (bool, int64, []*model.Message, error) {
	userMinSeq, err := s.seqUserCache.GetSeqUserMinSeq(ctx, userID, conversationID)
	if err != nil {
		return false, 0, nil, err
	}
	userMaxSeq, err := s.seqUserCache.GetSeqUserMaxSeq(ctx, userID, conversationID)
	if err != nil {
		return false, 0, nil, err
	}
	minSeq, err := s.seqConvCache.GetMinSeq(ctx, conversationID)
	if err != nil {
		return false, 0, nil, err
	}
	maxSeq, err := s.seqConvCache.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return false, 0, nil, err
	}
	if userMinSeq > minSeq {
		minSeq = userMinSeq
	}
	if userMaxSeq != 0 && userMaxSeq < maxSeq {
		maxSeq = userMaxSeq
	}
	var validSeqs []int64
	var (
		isEnd  bool
		endSeq int64
	)
	for _, seq := range seqs {
		if seq >= minSeq && seq <= maxSeq {
			validSeqs = append(validSeqs, seq)
		} else if seq < minSeq && pullOrder == PullOrderDesc {
			isEnd = true
			endSeq = minSeq
		} else if seq > maxSeq && pullOrder == PullOrderAsc {
			isEnd = true
			endSeq = maxSeq
		}
	}
	if len(validSeqs) == 0 {
		return isEnd, endSeq, nil, nil
	}
	msgs, err := s.GetMessageBySeqs(ctx, conversationID, userID, validSeqs)
	if err != nil {
		return false, 0, nil, err
	}
	return isEnd, endSeq, msgs, nil
}

type GetLastMessageReq struct {
	UserID          int64    `json:"user_id,string"`
	ConversationIDs []string `json:"conversation_ids"`
}
type GetLastMessageResp struct {
	Messages map[string]*model.Message `json:"messages"`
}

func (s *MessageService) GetLastMessage(ctx context.Context, req GetLastMessageReq) (GetLastMessageResp, error) {
	lastMessages := make(map[string]*model.Message)
	for _, convID := range req.ConversationIDs {
		userMaxSeq, err := s.seqUserCache.GetSeqUserMaxSeq(ctx, req.UserID, convID)
		if err != nil {
			return GetLastMessageResp{}, err
		}
		maxSeq, err := s.seqConvCache.GetMaxSeq(ctx, convID)
		if err != nil {
			return GetLastMessageResp{}, err
		}
		if userMaxSeq != 0 && userMaxSeq < maxSeq {
			maxSeq = userMaxSeq
		}
		if maxSeq == 0 {
			continue
		}

		// Check MinSeq 清空历史记录
		userMinSeq, err := s.seqUserCache.GetSeqUserMinSeq(ctx, req.UserID, convID)
		if err != nil {
			return GetLastMessageResp{}, err
		}
		minSeq, err := s.seqConvCache.GetMinSeq(ctx, convID)
		if err != nil {
			return GetLastMessageResp{}, err
		}
		if userMinSeq > minSeq {
			minSeq = userMinSeq
		}
		if maxSeq < minSeq {
			continue
		}

		msgs, err := s.GetMessageBySeqs(ctx, convID, req.UserID, []int64{maxSeq})
		if err != nil {
			return GetLastMessageResp{}, err
		}
		if len(msgs) > 0 {
			lastMessages[convID] = msgs[0]
		}
	}
	return GetLastMessageResp{Messages: lastMessages}, nil
}

type GetConversationsHasReadAndMaxSeqReq struct {
	UserID          int64    `json:"user_id,string"`
	ConversationIDs []string `json:"conversation_ids"`
}

type Seqs struct {
	MaxSeq     int64 `json:"max_seq"`
	HasReadSeq int64 `json:"has_read_seq"`
	MaxSeqTime int64 `json:"max_seq_time"`
}

type GetConversationsHasReadAndMaxSeqResp struct {
	Seqs map[string]*Seqs `json:"seqs"`
}

func (s *MessageService) GetConversationsHasReadAndMaxSeq(ctx context.Context, req GetConversationsHasReadAndMaxSeqReq) (GetConversationsHasReadAndMaxSeqResp, error) {
	resp := GetConversationsHasReadAndMaxSeqResp{
		Seqs: make(map[string]*Seqs),
	}
	maxSeqs, err := s.seqConvCache.GetMaxSeqs(ctx, req.ConversationIDs)
	if err != nil {
		return GetConversationsHasReadAndMaxSeqResp{}, err
	}
	var conversations []model.Conversation
	if err := s.db.WithContext(ctx).Where("owner_id = ? AND conversation_id IN ?", req.UserID, req.ConversationIDs).Find(&conversations).Error; err != nil {
		return GetConversationsHasReadAndMaxSeqResp{}, err
	}
	convMap := make(map[string]model.Conversation)
	for _, c := range conversations {
		convMap[c.ConversationID] = c
	}
	for _, convID := range req.ConversationIDs {
		var hasReadSeq int64
		if c, ok := convMap[convID]; ok {
			hasReadSeq = c.ReadSeq
		}
		resp.Seqs[convID] = &Seqs{
			MaxSeq:     maxSeqs[convID],
			HasReadSeq: hasReadSeq,
		}
	}
	return resp, nil
}

// ===================== Initialization Functions =====================

type InitConversationReq struct {
	SenderID int64 `json:"sender_id,string"`
	ConvType int32 `json:"conv_type" binding:"required"`
	TargetID int64 `json:"target_id,string" binding:"required"`
}

// 发信息时调用，初始化会话记录
func (s *MessageService) InitConversation(ctx context.Context, req InitConversationReq) error {
	conversationID := GetConversationID(req.ConvType, req.SenderID, req.TargetID)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch req.ConvType {
		case constant.SingleChatType:
			// 确保双方都有会话记录
			conversations := []model.Conversation{
				{OwnerID: req.SenderID, ConversationID: conversationID},
				{OwnerID: req.TargetID, ConversationID: conversationID},
			}
			for _, conversation := range conversations {
				if err := tx.FirstOrCreate(&conversation).Error; err != nil {
					return err
				}
			}
		case constant.GroupChatType:
			// TODO: 获取所有群成员ID，为所有没有群组会话记录的成员创建会话记录

		default:
			return errors.New("invalid session type")
		}
		return nil
	})
}

// ===================== Helper Functions =====================

func GetConversationID(ConvType int32, userID int64, receiverID int64) string {
	switch ConvType {
	case constant.SingleChatType:
		if userID <= receiverID {
			return "single:" + strconv.FormatInt(userID, 10) + "_" + strconv.FormatInt(receiverID, 10)
		}
		return "single:" + strconv.FormatInt(receiverID, 10) + "_" + strconv.FormatInt(userID, 10)
	case constant.GroupChatType:
		return "group:" + strconv.FormatInt(receiverID, 10)
	default:
		return ""
	}
}
