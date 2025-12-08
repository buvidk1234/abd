package imrepo

import (
	"backend/internal/model"
	"backend/internal/pkg/cache/cachekey"
	"backend/internal/pkg/database"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestImRepo(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.6.130:6379",
		Password: "123456",
		DB:       0,
	})
	db := database.GetDB()
	db.AutoMigrate(model.SeqConversation{})

	imRepo := NewImRepo(db, rdb)
	// 手动初始化 seqConv，因为 NewImRepo 目前没有做这件事

	if imRepo == nil {
		t.Errorf("NewImRepo returned nil")
	}

	ctx := context.Background()
	conversationID := "test_conversation_batch_insert"

	// 1. 构造批量消息
	msgCount := 5
	msgs := make([]*model.Message, 0, msgCount)
	for i := 0; i < msgCount; i++ {
		msgs = append(msgs, &model.Message{
			ConversationID: conversationID,
			SenderID:       "user_test_001",
			Content:        ("hello world " + time.Now().String()),
			MsgType:        1,
			SendTime:       time.Now().UnixMilli(),
		})
	}

	// 2. 执行批量插入
	err := imRepo.BatchStoreMsgToRedis(ctx, conversationID, msgs)
	if err != nil {
		t.Fatalf("BatchStoreMsgToRedis failed: %v", err)
	}
	t.Logf("BatchStoreMsgToRedis success, stored %d messages", len(msgs))

	// 3. 验证并打印插入的数据
	for i, msg := range msgs {
		if msg.Seq <= 0 {
			t.Errorf("Message[%d] Seq invalid: %d", i, msg.Seq)
		}

		// 从 Redis 读取验证
		key := cachekey.GetMsgCacheKey(conversationID, msg.Seq)
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			t.Errorf("Failed to get msg from redis: key=%s, err=%v", key, err)
			continue
		}

		var storedMsg model.Message
		err = json.Unmarshal([]byte(val), &storedMsg)
		if err != nil {
			t.Errorf("Failed to unmarshal msg: %v", err)
			continue
		}

		t.Logf("Checked Msg[%d] Redis: Seq=%d, Content=%s", i, storedMsg.Seq, string(storedMsg.Content))
	}
}

func TestBatchStoreMsgToDB(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.6.130:6379",
		Password: "123456",
		DB:       0,
	})
	db := database.GetDB()
	db.AutoMigrate(&model.Message{})

	imRepo := NewImRepo(db, rdb)
	ctx := context.Background()
	conversationID := "test_conv_db_only"

	// 构造测试数据
	msgs := []*model.Message{
		{
			ID:             1001,
			ConversationID: conversationID,
			Seq:            101,
			SenderID:       "user_A",
			Content:        "msg 1 content",
			MsgType:        1,
			SendTime:       time.Now().UnixMilli(),
			CreateTime:     time.Now().UnixMilli(),
		},
		{
			ID:             1002,
			ConversationID: conversationID,
			Seq:            102,
			SenderID:       "user_A",
			Content:        "msg 2 content",
			MsgType:        1,
			SendTime:       time.Now().UnixMilli(),
			CreateTime:     time.Now().UnixMilli(),
		},
	}

	// 清理旧数据
	db.Where("conversation_id = ?", conversationID).Delete(&model.Message{})

	// 执行入库
	err := imRepo.BatchStoreMsgToDB(ctx, msgs)
	if err != nil {
		t.Fatalf("BatchStoreMsgToDB failed: %v", err)
	}
	t.Log("BatchStoreMsgToDB executed successfully")

	// 验证数据
	var count int64
	db.Model(&model.Message{}).Where("conversation_id = ?", conversationID).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 messages in DB, got %d", count)
	}

	var storedMsgs []model.Message
	db.Where("conversation_id = ?", conversationID).Order("seq asc").Find(&storedMsgs)

	for i, msg := range storedMsgs {
		expected := msgs[i]
		if msg.Seq != expected.Seq {
			t.Errorf("Msg[%d] Seq mismatch: expected %d, got %d", i, expected.Seq, msg.Seq)
		}
		if msg.Content != expected.Content {
			t.Errorf("Msg[%d] Content mismatch: expected %s, got %s", i, expected.Content, msg.Content)
		}
		t.Logf("Verified Msg: Seq=%d, Content=%s", msg.Seq, msg.Content)
	}
}
