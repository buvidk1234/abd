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
	resp, err := svc.PullMessageBySeqs(ctx, userID, req)
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
	resp2, err := svc.PullMessageBySeqs(ctx, userID, req2)
	if err != nil {
		t.Fatalf("Case 2 failed: %v", err)
	}
	if len(resp2.Msgs[convID].Msgs) != 5 {
		t.Errorf("Case 2: expected 5 msgs, got %d", len(resp2.Msgs[convID].Msgs))
	}

	// Case 3: Limit Num
	req3 := PullMessageBySeqsReq{
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
	resp3, err := svc.PullMessageBySeqs(ctx, userID, req3)
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
	resp4, err := svc.PullMessageBySeqs(ctx, userID, req4)
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

	resp, err := svc.GetMaxSeq(ctx, userID)

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
	resp2, err := svc.GetMaxSeq(ctx, userID)
	if err != nil {
		t.Fatalf("GetMaxSeq (2nd call) failed: %v", err)
	}
	if resp2.MaxSeqs[convID] != maxSeq {
		t.Errorf("expected max seq %d on 2nd call, got %d", maxSeq, resp2.MaxSeqs[convID])
	}
}

func TestGetSeqMessage(t *testing.T) {
	// 1. Setup
	db := setupTestDB(t)
	setupRedis(t)
	redis.RDB.FlushDB(context.Background())

	svc := NewMessageService(db)
	ctx := context.Background()

	userID := int64(2001)
	msgIDCounter := int64(1000)

	// Helper to create conversation and messages
	createConv := func(convID string, userMinSeq, userMaxSeq int64) {
		// Conversation
		db.Create(&model.Conversation{
			OwnerID:        userID,
			ConversationID: convID,
			ConvType:       1,
			MaxSeq:         20,
		})
		// SeqConversation (Global)
		db.Create(&model.SeqConversation{
			ID:      convID,
			SeqType: 1,
			MaxSeq:  20,
			MinSeq:  0,
		})
		// SeqUser
		db.Create(&model.SeqUser{
			UserID:         userID,
			ConversationID: convID,
			MinSeq:         userMinSeq,
			MaxSeq:         userMaxSeq,
		})
		// Messages
		for i := int64(1); i <= 20; i++ {
			msgIDCounter++
			db.Create(&model.Message{
				ID:             msgIDCounter, // Unique ID
				ConversationID: convID,
				Seq:            i,
				SenderID:       userID,
				Content:        "msg-" + strconv.FormatInt(i, 10),
				MsgType:        1,
			})
		}
	}

	// Scenario 1: Normal (No limits)
	convNormal := "single:2001_2002"
	createConv(convNormal, 0, 0)

	// Scenario 2: MinSeq Limit (History Cleared up to 5)
	convMin := "single:2001_2003"
	createConv(convMin, 5, 0)

	// Scenario 3: MaxSeq Limit (Left Group at 15)
	convMax := "group:2004"
	createConv(convMax, 0, 15)

	// 2. Test Request
	req := GetSeqMessageReq{
		UserID: userID,
		Conversations: []*ConversationSeqs{
			{
				ConversationID: convNormal,
				Seqs:           []int64{1, 2, 3},
			},
			{
				ConversationID: convMin,
				Seqs:           []int64{3, 4, 5, 6}, // 3,4 should be filtered
			},
			{
				ConversationID: convMax,
				Seqs:           []int64{14, 15, 16, 17}, // 16,17 should be filtered
			},
		},
		Order: PullOrderAsc,
	}

	// 3. Execute
	resp, err := svc.GetSeqMessage(ctx, req)
	if err != nil {
		t.Fatalf("GetSeqMessage failed: %v", err)
	}

	// 4. Assertions

	// Assert Normal
	if msgs := resp.Msgs[convNormal]; len(msgs.Msgs) != 3 {
		t.Errorf("Normal: expected 3 msgs, got %d", len(msgs.Msgs))
	}

	// Assert MinSeq
	if msgs := resp.Msgs[convMin]; len(msgs.Msgs) != 2 { // Only 5, 6
		t.Errorf("MinSeq: expected 2 msgs (5,6), got %d", len(msgs.Msgs))
	} else {
		if msgs.Msgs[0].Seq != 5 {
			t.Errorf("MinSeq: expected first msg seq 5, got %d", msgs.Msgs[0].Seq)
		}
	}

	// Assert MaxSeq
	if msgs := resp.Msgs[convMax]; len(msgs.Msgs) != 2 { // Only 14, 15
		t.Errorf("MaxSeq: expected 2 msgs (14,15), got %d", len(msgs.Msgs))
	} else {
		if msgs.IsEnd != true {
			t.Errorf("MaxSeq: expected IsEnd true (requested > maxSeq in Asc)")
		}
		if msgs.EndSeq != 15 {
			t.Errorf("MaxSeq: expected EndSeq 15, got %d", msgs.EndSeq)
		}
	}

	// 5. Test Descending Order for MinSeq boundary
	reqDesc := GetSeqMessageReq{
		UserID: userID,
		Conversations: []*ConversationSeqs{
			{
				ConversationID: convMin,
				Seqs:           []int64{6, 5, 4, 3}, // 4,3 filtered
			},
		},
		Order: PullOrderDesc,
	}
	respDesc, err := svc.GetSeqMessage(ctx, reqDesc)
	if err != nil {
		t.Fatalf("GetSeqMessage Desc failed: %v", err)
	}
	if msgs := respDesc.Msgs[convMin]; len(msgs.Msgs) != 2 {
		t.Errorf("Desc MinSeq: expected 2 msgs, got %d", len(msgs.Msgs))
	} else {
		if msgs.IsEnd != true {
			t.Errorf("Desc MinSeq: expected IsEnd true")
		}
		if msgs.EndSeq != 5 { // MinSeq is 5
			t.Errorf("Desc MinSeq: expected EndSeq 5, got %d", msgs.EndSeq)
		}
	}
}

func TestGetLastMessage(t *testing.T) {
	// 1. Setup
	db := setupTestDB(t)
	setupRedis(t)
	redis.RDB.FlushDB(context.Background())

	svc := NewMessageService(db)
	ctx := context.Background()

	userID := int64(3001)
	msgIDCounter := int64(3000)

	// Helper
	createConv := func(convID string, userMinSeq, userMaxSeq int64) {
		db.Create(&model.Conversation{
			OwnerID:        userID,
			ConversationID: convID,
			ConvType:       1,
			MaxSeq:         20,
		})
		db.Create(&model.SeqConversation{
			ID:      convID,
			SeqType: 1,
			MaxSeq:  20,
			MinSeq:  0,
		})
		db.Create(&model.SeqUser{
			UserID:         userID,
			ConversationID: convID,
			MinSeq:         userMinSeq,
			MaxSeq:         userMaxSeq,
		})
		for i := int64(1); i <= 20; i++ {
			msgIDCounter++
			db.Create(&model.Message{
				ID:             msgIDCounter,
				ConversationID: convID,
				Seq:            i,
				SenderID:       userID,
				Content:        "msg-" + strconv.FormatInt(i, 10),
				MsgType:        1,
			})
		}
	}

	// Scenario 1: Normal
	convNormal := "single:3001_3002"
	createConv(convNormal, 0, 0)

	// Scenario 2: Cleared History (MinSeq > MaxSeq)
	// User cleared history when MaxSeq was 20. So MinSeq becomes 21.
	convCleared := "single:3001_3003"
	createConv(convCleared, 21, 0)

	// Scenario 3: User Left (MaxSeq Limit)
	// Global Max 20. User Max 15. Should see 15.
	convLeft := "group:3004"
	createConv(convLeft, 0, 15)

	// Scenario 4: Partial History (MinSeq set)
	// MinSeq = 18. Should see 20.
	convPartial := "single:3001_3005"
	createConv(convPartial, 18, 0)

	req := GetLastMessageReq{
		UserID:          userID,
		ConversationIDs: []string{convNormal, convCleared, convLeft, convPartial},
	}

	resp, err := svc.GetLastMessage(ctx, req)
	if err != nil {
		t.Fatalf("GetLastMessage failed: %v", err)
	}

	// Assert Normal
	if msg, ok := resp.Messages[convNormal]; !ok || msg.Seq != 20 {
		t.Errorf("Normal: expected seq 20, got %v", msg)
	}

	// Assert Cleared
	if msg, ok := resp.Messages[convCleared]; ok {
		t.Errorf("Cleared: expected no message, got seq %d", msg.Seq)
	}

	// Assert Left
	if msg, ok := resp.Messages[convLeft]; !ok || msg.Seq != 15 {
		t.Errorf("Left: expected seq 15, got %v", msg)
	}

	// Assert Partial
	if msg, ok := resp.Messages[convPartial]; !ok || msg.Seq != 20 {
		t.Errorf("Partial: expected seq 20, got %v", msg)
	}
}

func TestGetConversationsHasReadAndMaxSeq(t *testing.T) {
	// 1. Setup
	db := setupTestDB(t)
	setupRedis(t)
	redis.RDB.FlushDB(context.Background())
	svc := NewMessageService(db)
	ctx := context.Background()

	userID := int64(2001)
	convID1 := "single:2001_2002"
	convID2 := "group:3001"

	// 2. Prepare Data
	// 2.1 Conversations
	// conv1: ReadSeq=5
	if err := db.Create(&model.Conversation{
		OwnerID:        userID,
		ConversationID: convID1,
		ConvType:       1,
		ReadSeq:        5,
	}).Error; err != nil {
		t.Fatalf("failed to create conv1: %v", err)
	}
	// conv2: ReadSeq=10
	if err := db.Create(&model.Conversation{
		OwnerID:        userID,
		ConversationID: convID2,
		ConvType:       2,
		ReadSeq:        10,
	}).Error; err != nil {
		t.Fatalf("failed to create conv2: %v", err)
	}

	// 2.2 SeqConversation (MaxSeq in Redis/DB)
	// conv1: MaxSeq=20
	if err := db.Create(&model.SeqConversation{
		ID:     convID1,
		MaxSeq: 20,
	}).Error; err != nil {
		t.Fatalf("failed to create seq conv1: %v", err)
	}
	// conv2: MaxSeq=30
	if err := db.Create(&model.SeqConversation{
		ID:     convID2,
		MaxSeq: 30,
	}).Error; err != nil {
		t.Fatalf("failed to create seq conv2: %v", err)
	}

	// 3. Test Case
	req := GetConversationsHasReadAndMaxSeqReq{
		UserID:          userID,
		ConversationIDs: []string{convID1, convID2, "non_existent_conv"},
	}

	resp, err := svc.GetConversationsHasReadAndMaxSeq(ctx, req)
	if err != nil {
		t.Fatalf("GetConversationsHasReadAndMaxSeq failed: %v", err)
	}

	// 4. Verify Results
	// Verify conv1
	if seqs, ok := resp.Seqs[convID1]; !ok {
		t.Errorf("conv1 not found in response")
	} else {
		if seqs.HasReadSeq != 5 {
			t.Errorf("conv1 HasReadSeq expected 5, got %d", seqs.HasReadSeq)
		}
		if seqs.MaxSeq != 20 {
			t.Errorf("conv1 MaxSeq expected 20, got %d", seqs.MaxSeq)
		}
	}

	// Verify conv2
	if seqs, ok := resp.Seqs[convID2]; !ok {
		t.Errorf("conv2 not found in response")
	} else {
		if seqs.HasReadSeq != 10 {
			t.Errorf("conv2 HasReadSeq expected 10, got %d", seqs.HasReadSeq)
		}
		if seqs.MaxSeq != 30 {
			t.Errorf("conv2 MaxSeq expected 30, got %d", seqs.MaxSeq)
		}
	}

	// Verify non_existent_conv
	if seqs, ok := resp.Seqs["non_existent_conv"]; !ok {
		t.Errorf("non_existent_conv not found in response")
	} else {
		if seqs.HasReadSeq != 0 {
			t.Errorf("non_existent_conv HasReadSeq expected 0, got %d", seqs.HasReadSeq)
		}
		// MaxSeq might be 0 if not found
		if seqs.MaxSeq != 0 {
			t.Errorf("non_existent_conv MaxSeq expected 0, got %d", seqs.MaxSeq)
		}
	}
}
