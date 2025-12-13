package service

import (
	"backend/internal/model"
	"backend/internal/pkg/cache/cachekey"
	"backend/internal/pkg/cache/redis"
	"context"
	"strconv"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	// Clear old data
	if err := db.Migrator().DropTable(&model.Conversation{}, &model.SeqConversation{}, &model.Message{}, &model.SeqUser{}); err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}
	err = db.AutoMigrate(&model.Conversation{}, &model.SeqConversation{}, &model.Message{}, &model.SeqUser{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func setupRedis(t *testing.T) {
	cfg := redis.Config{
		Addr:         "192.168.6.130:6379",
		Password:     "123456",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
	}
	redis.Init(cfg)
}

func TestPullMessageBySeqs(t *testing.T) {
	// 1. Setup
	db := setupTestDB(t)
	setupRedis(t)

	redis.RDB.FlushDB(context.Background())
	svc := NewMessageService(db)
	ctx := context.Background()

	userID := int64(1001)
	convID := "single:1001_1002"
	// 2. Prepare Data
	// 2.1 Conversation
	conv := model.Conversation{
		OwnerID:        userID,
		ConversationID: convID,
		ConvType:       1,
		MaxSeq:         10, // userMaxSeq
	}
	if err := db.Create(&conv).Error; err != nil {
		t.Fatalf("failed to create conversation: %v", err)
	}

	// 2.2 SeqConversation (Global Min/Max for conversation)
	seqConv := model.SeqConversation{
		ID:      convID,
		SeqType: 1,
		MaxSeq:  100,
		MinSeq:  0,
	}
	if err := db.Create(&seqConv).Error; err != nil {
		t.Fatalf("failed to create seq conversation: %v", err)
	}

	// 2.3 SeqUser (User specific MinSeq)
	seqUser := model.SeqUser{
		UserID:         userID,
		ConversationID: convID,
		MinSeq:         0,
		MaxSeq:         100,
	}
	if err := db.Create(&seqUser).Error; err != nil {
		t.Fatalf("failed to create seq user: %v", err)
	}

	// 2.4 Messages
	for i := int64(1); i <= 10; i++ {
		msg := model.Message{
			ID:             i,
			ConversationID: convID,
			Seq:            i,
			SenderID:       userID,
			Content:        "msg",
			MsgType:        1,
		}
		if err := db.Create(&msg).Error; err != nil {
			t.Fatalf("failed to create message: %v", err)
		}
	}

	// 3. Test Cases

	// Case 1: Pull Ascending 1-5
	req := PullMessageBySeqsReq{
		UserID: userID,
		SeqRanges: []*SeqRange{
			{
				ConversationID: convID,
				Begin:          1,
				End:            5,
				Num:            10,
			},
		},
		Order: PullOrderAsc,
	}
	resp, err := svc.PullMessageBySeqs(ctx, req)
	if err != nil {
		t.Fatalf("Case 1 failed: %v", err)
	}
	if len(resp.Msgs[convID].Msgs) != 5 {
		t.Errorf("Case 1: expected 5 msgs, got %d", len(resp.Msgs[convID].Msgs))
	}
	if resp.Msgs[convID].IsEnd != false {
		t.Errorf("Case 1: expected IsEnd false")
	}

	// Case 2: Pull Descending 10-6
	req2 := PullMessageBySeqsReq{
		UserID: userID,
		SeqRanges: []*SeqRange{
			{
				ConversationID: convID,
				Begin:          6,
				End:            10,
				Num:            10,
			},
		},
		Order: PullOrderDesc,
	}
	resp2, err := svc.PullMessageBySeqs(ctx, req2)
	if err != nil {
		t.Fatalf("Case 2 failed: %v", err)
	}
	if len(resp2.Msgs[convID].Msgs) != 5 {
		t.Errorf("Case 2: expected 5 msgs, got %d", len(resp2.Msgs[convID].Msgs))
	}

	// Case 3: Limit Num
	req3 := PullMessageBySeqsReq{
		UserID: userID,
		SeqRanges: []*SeqRange{
			{
				ConversationID: convID,
				Begin:          1,
				End:            10,
				Num:            3,
			},
		},
		Order: PullOrderAsc,
	}
	resp3, err := svc.PullMessageBySeqs(ctx, req3)
	if err != nil {
		t.Fatalf("Case 3 failed: %v", err)
	}
	if len(resp3.Msgs[convID].Msgs) != 3 {
		t.Errorf("Case 3: expected 3 msgs, got %d", len(resp3.Msgs[convID].Msgs))
	}

	// Case 4: MaxSeq Limit (User Left Group)
	// Simulate user leaving group by setting MaxSeq to 5.
	// This means the user should not be able to pull messages with Seq > 5.
	if err := db.Model(&model.Conversation{}).Where("owner_id = ? AND conversation_id = ?", userID, convID).Update("max_seq", 5).Error; err != nil {
		t.Fatalf("failed to update max_seq: %v", err)
	}

	// Delete cache to ensure the updated MaxSeq is fetched from DB
	key := cachekey.GetConversationKey(strconv.FormatInt(userID, 10), convID)
	if err := redis.RDB.Del(ctx, key).Err(); err != nil {
		t.Fatalf("failed to delete cache: %v", err)
	}

	var updatedConv model.Conversation
	if err := db.Where("owner_id = ? AND conversation_id = ?", userID, convID).First(&updatedConv).Error; err != nil {
		t.Fatalf("failed to fetch updated conversation: %v", err)
	}
	t.Logf("updatedConv: %v", updatedConv)
	req4 := PullMessageBySeqsReq{
		UserID: userID,
		SeqRanges: []*SeqRange{
			{
				ConversationID: convID,
				Begin:          6,
				End:            10,
				Num:            5,
			},
		},
		Order: PullOrderAsc,
	}
	resp4, err := svc.PullMessageBySeqs(ctx, req4)
	if err != nil {
		t.Fatalf("Case 4 failed: %v", err)
	}
	// Expecting no messages because the requested range (6-10) is entirely beyond SyncSeq (5).
	if val, ok := resp4.Msgs[convID]; ok {
		t.Errorf("Case 4: expected no msgs, got %d", len(val.Msgs))
	}
}

func TestGetMaxSeq(t *testing.T) {
	// 1. Setup
	db := setupTestDB(t)
	setupRedis(t)

	svc := NewMessageService(db)
	ctx := context.Background()

	userID := int64(1001)
	convID := "single:1001_1002"
	maxSeq := int64(100)

	// 2. Prepare Data
	// Create Conversation
	conv := model.Conversation{
		OwnerID:        1001,
		ConversationID: convID,
		ConvType:       1,
		Status:         1,
	}
	if err := db.Create(&conv).Error; err != nil {
		t.Fatalf("failed to create conversation: %v", err)
	}

	// Create SeqConversation
	seqConv := model.SeqConversation{
		ID:      convID,
		SeqType: 1,
		MaxSeq:  maxSeq,
	}
	if err := db.Create(&seqConv).Error; err != nil {
		t.Fatalf("failed to create seq conversation: %v", err)
	}

	// Clear Redis Cache to ensure we hit DB or fresh cache logic
	// Note: In a real environment, be careful clearing keys.
	// Here we assume test environment.
	// We can delete specific keys used by the test.
	// Key for conversation IDs: "conversation_ids:1001" (assuming cachekey logic)
	// Key for max seq: "seq_conv:single:1001_1002" (assuming cachekey logic)
	// Since we don't have easy access to cachekey package functions here without importing,
	// we rely on the service to work correctly.
	// Ideally we should import cachekey package.

	// 3. Execute
	req := GetMaxSeqReq{
		UserID: userID,
	}
	resp, err := svc.GetMaxSeq(ctx, req)

	// 4. Assert
	if err != nil {
		t.Fatalf("GetMaxSeq failed: %v", err)
	}

	if len(resp.MaxSeqs) != 1 {
		t.Errorf("expected 1 max seq, got %d", len(resp.MaxSeqs))
	}

	if seq, ok := resp.MaxSeqs[convID]; !ok {
		t.Errorf("expected conversation %s in response", convID)
	} else {
		if seq != maxSeq {
			t.Errorf("expected max seq %d, got %d", maxSeq, seq)
		}
	}

	// Test Cache Hit (Optional: Run again and check logs or coverage, but hard to assert without mocking)
	resp2, err := svc.GetMaxSeq(ctx, req)
	if err != nil {
		t.Fatalf("GetMaxSeq (2nd call) failed: %v", err)
	}
	if resp2.MaxSeqs[convID] != maxSeq {
		t.Errorf("expected max seq %d on 2nd call, got %d", maxSeq, resp2.MaxSeqs[convID])
	}
}
