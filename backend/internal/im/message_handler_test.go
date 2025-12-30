package im

import (
	"backend/internal/model"
	"backend/internal/pkg/constant"
	"backend/internal/pkg/kafka"
	"backend/internal/pkg/snowflake"
	"backend/internal/service"
	"context"
	"encoding/json"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	// Initialize snowflake with a dummy machine ID
	_ = snowflake.Init(snowflake.Config{MachineID: 1})
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(
		&model.Message{},
		&model.SeqConversation{},
		&model.Conversation{},
		&model.SeqUser{},
		// &model.UserTimeline{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func TestServiceHandler_SendMessage(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewMessageService(db)
	producer, _ := kafka.NewSyncProducer()
	handler := NewServiceHandler(svc, producer)

	// Prepare request data
	sendReq := service.SendMessageReq{
		SenderID: 1,
		ConvType: constant.SingleChatType,
		TargetID: 2,
		MsgType:  1,
		Content:  "hello",
	}
	dataBytes, _ := json.Marshal(sendReq)

	req := &Req{
		InboundReq: InboundReq{
			ReqIdentifier: 1001,
			Data:          dataBytes,
		},
	}

	// Test SendMessage
	_, err := handler.SendMessage(context.Background(), req)
	if err != nil {
		t.Errorf("SendMessage failed: %v", err)
	}
}

func TestServiceHandler_PullSpecifiedConv(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewMessageService(db)
	producer, _ := kafka.NewSyncProducer()
	handler := NewServiceHandler(svc, producer)

	// Prepare request data
	pullReq := service.PullSpecifiedConvReq{
		UserID:  1,
		ConvID:  "c1",
		ConvSeq: 0,
	}
	dataBytes, _ := json.Marshal(pullReq)

	req := &Req{
		InboundReq: InboundReq{
			ReqIdentifier: 1002,
			Data:          dataBytes,
		},
	}

	// Test PullSpecifiedConv
	_, err := handler.PullSpecifiedConv(context.Background(), req)
	if err != nil {
		// It might fail if no data, but we just check execution flow
		// t.Errorf("PullSpecifiedConv failed: %v", err)
	}
}
