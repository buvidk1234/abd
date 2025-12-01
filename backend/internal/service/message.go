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
	SenderID    string `json:"sender_id" binding:"required"`
	SessionType int32  `json:"session_type" binding:"required"`
	ReceiverID  string `json:"receiver_id" binding:"required"`
	ContentType int32  `json:"content_type" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

func (s *MessageService) SendMessage(ctx context.Context, req SendMessageReq) error {
	conversationID := s.getConversationID(req.SessionType, req.SenderID, req.ReceiverID)

	// TODO: 创建好友时或加入群聊时调用，初始化会话记录
	s.InitConversationForUser(ctx, req)

	var seqConversation model.SeqConversation
	// perform read->increment->update and message create in a single transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// lock the seq row to avoid concurrent increments
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("conversation_id = ?", conversationID).First(&seqConversation).Error; err != nil {
			return err
		}

		seqConversation.MaxSeq = seqConversation.MaxSeq + 1

		if err := tx.Model(&model.SeqConversation{}).Where("conversation_id = ?", conversationID).Update("max_seq", seqConversation.MaxSeq).Error; err != nil {
			return err
		}

		msg := model.Message{
			ConversationID: conversationID,
			SenderID:       req.SenderID,
			ContentType:    int32(req.ContentType),
			Content:        req.Content,
			Seq:            seqConversation.MaxSeq,
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

	return nil
}

type PullMsgResp struct {
	Messages []model.Message `json:"messages"`
	ReadSeq  int64           `json:"read_seq,string"`
	MaxSeq   int64           `json:"max_seq,string"`
}

func (s *MessageService) PullSpecifiedMsg(ctx context.Context, userAID string, sesstionType int32, oppositeID string) (PullMsgResp, error) {
	conversationID := s.getConversationID(sesstionType, userAID, oppositeID)
	var maxSeq int64
	s.db.WithContext(ctx).Model(&model.SeqConversation{}).Where("conversation_id = ?", conversationID).Select("max_seq").Scan(&maxSeq)
	var imaxSeq int64
	s.db.WithContext(ctx).Model(&model.Conversation{}).Where("conversation_id = ? AND user_id = ?", conversationID, userAID).Select("read_seq").Scan(&imaxSeq)
	if imaxSeq >= maxSeq {
		return PullMsgResp{
			Messages: []model.Message{},
			ReadSeq:  imaxSeq,
			MaxSeq:   maxSeq,
		}, nil
	}
	var messages []model.Message
	s.db.WithContext(ctx).Find(&messages, "conversation_id = ? and seq > ? and seq <= ?", conversationID, imaxSeq, maxSeq).Order("seq ASC")
	return PullMsgResp{
		Messages: messages,
		ReadSeq:  imaxSeq,
		MaxSeq:   maxSeq,
	}, nil
}

type PullAllMsgResp struct {
	PullMsgs map[string]PullMsgResp `json:"pull_msgs"`
}

func (s *MessageService) PullAllMsg(ctx context.Context, userID string) (PullAllMsgResp, error) {
	var conversations []model.Conversation
	s.db.WithContext(ctx).Find(&conversations, "user_id = ?", userID)
	var conversationIDs []string
	for _, conv := range conversations {
		conversationIDs = append(conversationIDs, conv.ConversationID)
	}
	var seqConversations []model.SeqConversation
	s.db.WithContext(ctx).Find(&seqConversations, "conversation_id IN ?", conversationIDs)
	seqMap := make(map[string]int64)
	for _, seqConv := range seqConversations {
		seqMap[seqConv.ConversationID] = seqConv.MaxSeq
	}

	pullMsgs := make(map[string]PullMsgResp)
	for _, conv := range conversations {
		maxSeq := seqMap[conv.ConversationID]
		if conv.ReadSeq < maxSeq {
			var messages []model.Message
			s.db.WithContext(ctx).Find(&messages, "conversation_id = ? and seq > ? and seq <= ?", conv.ConversationID, conv.ReadSeq, maxSeq).Order("seq ASC")
			pullMsgs[conv.ConversationID] = PullMsgResp{
				Messages: messages,
				ReadSeq:  conv.ReadSeq,
				MaxSeq:   maxSeq,
			}
		}
	}

	// Convert pullMsgs map to a sorted array based on the last message's sequence number
	// var sortedPullMsgs []PullMsgResp
	// for _, pullMsg := range pullMsgs {
	// 	sortedPullMsgs = append(sortedPullMsgs, pullMsg)
	// }
	// sort.Slice(sortedPullMsgs, func(i, j int) bool {
	// 	return sortedPullMsgs[i].MaxSeq < sortedPullMsgs[j].MaxSeq
	// })

	return PullAllMsgResp{PullMsgs: pullMsgs}, nil
}

type MarkMsgsAsReadReq struct {
	UserID      string `json:"user_id" binding:"required"`
	SessionType int32  `json:"session_type,string" binding:"required"`
	OppositeID  string `json:"opposite_id" binding:"required"`
	ReadSeq     int64  `json:"read_seq,string" binding:"required"`
}

func (s *MessageService) MarkMsgsAsRead(ctx context.Context, req MarkMsgsAsReadReq) error {
	conversationID := s.getConversationID(req.SessionType, req.UserID, req.OppositeID)
	return s.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, req.UserID).
		Update("read_seq", req.ReadSeq).Error
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
	conversationID := s.getConversationID(req.SessionType, req.SenderID, req.ReceiverID)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch req.SessionType {
		case constant.SingleChatType:
			// 确保双方都有会话记录
			conversations := []model.Conversation{
				{UserID: req.SenderID, ConversationID: conversationID},
				{UserID: req.ReceiverID, ConversationID: conversationID},
			}
			for _, conversation := range conversations {
				if err := tx.FirstOrCreate(&conversation).Error; err != nil {
					return err
				}
			}
			seqConversation := model.SeqConversation{
				ConversationID: conversationID,
				MinSeq:         0,
				MaxSeq:         0,
			}
			if err := tx.FirstOrCreate(&seqConversation, "conversation_id = ?", conversationID).Error; err != nil {
				return err
			}
		case constant.GroupChatType:
			// 确保群组会话记录存在
			var seqConversation model.SeqConversation
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("conversation_id = ?", conversationID).First(&seqConversation).Error; err != nil {
				return err
			}
			conversation := model.Conversation{
				UserID:         req.SenderID,
				ConversationID: conversationID,
				MinSeq:         seqConversation.MaxSeq,
				MaxSeq:         seqConversation.MaxSeq,
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
			ConversationID: conversationID,
			MinSeq:         0,
			MaxSeq:         0,
		}
		if err := tx.FirstOrCreate(&seqConversation, "conversation_id = ?", conversationID).Error; err != nil {
			return err
		}

		var conversations []model.Conversation
		for _, memberID := range memberIDs {
			conversations = append(conversations, model.Conversation{
				UserID:         memberID,
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

func (s *MessageService) getConversationID(sessionType int32, userID string, receiverID string) string {
	switch sessionType {
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
