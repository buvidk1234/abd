package service

import (
	"backend/internal/model"
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
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

type SendMessageReq struct {
	SenderID string `json:"sender_id"`
	ConvType int32  `json:"conv_type" binding:"required"`
	TargetID string `json:"target_id" binding:"required"`
	MsgType  int32  `json:"msg_type" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func (s *MessageService) SendMessage(ctx context.Context, req SendMessageReq) error {
	conversationID := GetConversationID(req.ConvType, req.SenderID, req.TargetID)

	// TODO: 创建好友时或加入群聊时调用，初始化会话记录
	s.InitConversationForUser(ctx, req)

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
	UserID  string `json:"user_id"`
	ConvID  string `json:"conv_id"`
	ConvSeq int32  `json:"conv_seq,string"`
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
	UserID  string `json:"user_id"`
	UserSeq int64  `json:"user_seq,string"`
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

func (s *MessageService) DeleteConversation(ctx context.Context, userID string, conversationID string) error {
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

// ===================== Initialization Functions =====================

// TODO: 成为好友，或加入群组时调用，初始化会话记录
func (s *MessageService) InitConversationForUser(ctx context.Context, req SendMessageReq) error {
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
			seqConversation := model.SeqConversation{
				ID:     conversationID,
				MaxSeq: 0,
			}
			if err := tx.FirstOrCreate(&seqConversation, "id = ?", conversationID).Error; err != nil {
				return err
			}
		case constant.GroupChatType:
			// 确保群组会话记录存在
			var seqConversation model.SeqConversation
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", conversationID).First(&seqConversation).Error; err != nil {
				return err
			}
			conversation := model.Conversation{
				OwnerID:        req.SenderID,
				ConversationID: conversationID,
				MinSeq:         seqConversation.MaxSeq,
				SyncSeq:        seqConversation.MaxSeq,
			}
			if err := tx.FirstOrCreate(&conversation).Error; err != nil {
				return err
			}
		default:
			return errors.New("invalid session type")
		}
		return nil
	})
}

// TODO: 创建群组时调用，初始化群组会话记录
func (s *MessageService) InitConversationForCreateGroup(ctx context.Context, groupID string, memberIDs []string) error {
	conversationID := GetConversationID(constant.GroupChatType, "", groupID)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		seqConversation := model.SeqConversation{
			ID:     conversationID,
			MaxSeq: 0,
		}
		if err := tx.FirstOrCreate(&seqConversation, "conversation_id = ?", conversationID).Error; err != nil {
			return err
		}

		var conversations []model.Conversation
		for _, memberID := range memberIDs {
			conversations = append(conversations, model.Conversation{
				OwnerID:        memberID,
				ConversationID: conversationID,
			})
		}
		if len(conversations) > 0 {
			if err := tx.Create(&conversations).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ===================== Helper Functions =====================

func GetConversationID(ConvType int32, userID string, receiverID string) string {
	switch ConvType {
	case constant.SingleChatType:
		v1, _ := strconv.ParseInt(userID, 10, 64)
		v2, _ := strconv.ParseInt(receiverID, 10, 64)
		if v1 <= v2 {
			return "single:" + userID + "_" + receiverID
		}
		return "single:" + receiverID + "_" + userID
	case constant.GroupChatType:
		return "group:" + receiverID
	default:
		return ""
	}
}
