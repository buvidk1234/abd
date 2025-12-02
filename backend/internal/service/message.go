package service

import (
	"backend/internal/model"
	"backend/internal/pkg/constant"
	"context"
	"errors"
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
	SenderID string `json:"sender_id" binding:"required"`
	ConvType int32  `json:"conv_type" binding:"required"`
	TargetID string `json:"target_id" binding:"required"`
	MsgType  int32  `json:"msg_type" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func (s *MessageService) SendMessage(ctx context.Context, req SendMessageReq) error {
	conversationID := s.getConversationID(req.ConvType, req.SenderID, req.TargetID)

	// TODO: 创建好友时或加入群聊时调用，初始化会话记录
	s.InitConversationForUser(ctx, req)

	var seqConversation model.SeqConversation
	// perform read->increment->update and message create in a single transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// lock the seq row to avoid concurrent increments
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", conversationID).First(&seqConversation).Error; err != nil {
			return err
		}

		seqConversation.MaxSeq = seqConversation.MaxSeq + 1

		if err := tx.Model(&model.SeqConversation{}).Where("id = ?", conversationID).Update("max_seq", seqConversation.MaxSeq).Error; err != nil {
			return err
		}

		msg := model.Message{
			ConversationID: conversationID,
			Seq:            seqConversation.MaxSeq,
			SenderID:       req.SenderID,
			MsgType:        req.MsgType,
			Content:        req.Content,
		}

		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		return nil
	})

	// if transaction returned record not found, normalize error for caller if needed
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	// 更新timeline
	switch req.ConvType {
	case constant.SingleChatType:
		// 更新双方会话的 max_seq
		var seqUser model.SeqUser
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

			seqUser = model.SeqUser{
				ID:     req.TargetID,
				MaxSeq: 0,
			}
			tx.FirstOrCreate(&seqUser, "id = ?", req.TargetID)

			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.TargetID).First(&seqUser).Error; err != nil {
				return err
			}

			seqUser.MaxSeq = seqUser.MaxSeq + 1
			if err := tx.Model(&model.SeqUser{}).Where("id = ?", req.TargetID).Update("max_seq", seqUser.MaxSeq).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		userTimeline := model.UserTimeline{
			OwnerID:        req.TargetID,
			Seq:            seqUser.MaxSeq,
			ConversationID: conversationID,
			MsgID:          0, // TODO: 关联消息ID
			RefMsgSeq:      seqConversation.MaxSeq,
			MsgType:        req.MsgType,
			SenderID:       req.SenderID,
			Snapshot:       req.Content, // TODO: 生成摘要
		}
		if err := s.db.WithContext(ctx).Create(&userTimeline).Error; err != nil {
			return err
		}
	case constant.GroupChatType:
		// TODO: 群组消息更新群成员的timeline
	default:
		return errors.New("invalid session type")
	}
	return nil
}

type PullSpecifiedConvReq struct {
	UserID  string `json:"user_id" binding:"required"`
	ConvID  string `json:"conv_id" binding:"required"`
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
	UserID  string `json:"user_id" binding:"required"`
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
	conversationID := s.getConversationID(req.ConvType, req.SenderID, req.TargetID)
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
	conversationID := s.getConversationID(constant.GroupChatType, "", groupID)
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

func (s *MessageService) getConversationID(ConvType int32, userID string, receiverID string) string {
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
