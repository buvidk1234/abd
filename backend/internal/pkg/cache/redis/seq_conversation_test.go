package redis

import (
	"backend/internal/model"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSeqConversation(t *testing.T) {
	RDB = redis.NewClient(&redis.Options{
		Addr:         "192.168.6.130:6379",
		Password:     "123456",
		DB:           0,
		PoolSize:     5,
		MinIdleConns: 2,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  3 * time.Second,
	})
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	seqGenerator := NewSeqConversationCacheRedis(db, RDB)
	RDB.FlushDB(context.Background())
	db.Migrator().DropTable(&model.SeqConversation{})
	db.AutoMigrate(model.SeqConversation{})
	for i := 0; i < 10; i++ {
		var seqConv model.SeqConversation

		if err := db.Find(&seqConv, "id = ?", "conversation_01").Error; err != nil {
			t.Fatalf("DB Find SeqConversation failed: %v", err)
		}
		if seqConv.ID == "" {
			t.Logf("%v not exists", "conversation_01")
		}

		seq, err := seqGenerator.Malloc(context.Background(), "conversation_01", 5)
		if err != nil {
			t.Fatalf("SeqConversation Malloc failed: %v", err)
		}
		t.Logf("SeqConversation Malloc success: %d", seq)
	}
	var seqConv model.SeqConversation
	db.Find(&seqConv, "id = ?", "conversation_01")
	t.Logf("Final SeqConversation in DB: %+v", seqConv)
	if seqConv.MaxSeq != 55 {
		t.Fatalf("Final SeqConversation MaxSeq expected 50, got %d", seqConv.MaxSeq)
	}
}

func TestGetMaxSeqs(t *testing.T) {
	// 1. 初始化 Redis 和 SQLite
	rdb := redis.NewClient(&redis.Options{
		Addr:         "192.168.6.130:6379",
		Password:     "123456",
		DB:           0,
		PoolSize:     5,
		MinIdleConns: 2,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  3 * time.Second,
	})
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.SeqConversation{})

	// 清理环境
	rdb.FlushDB(context.Background())
	db.Exec("DELETE FROM seq_conversations")

	seqCache := NewSeqConversationCacheRedis(db, rdb)
	ctx := context.Background()

	// 准备测试数据
	convID1 := "conv_001"
	convID2 := "conv_002"
	convID3 := "conv_003"

	// 场景 1: 缓存中存在 (先 Malloc 产生数据)
	// conv_001: Malloc 10 -> DB=10+50=60, Redis CURR=10
	if _, err := seqCache.Malloc(ctx, convID1, 10); err != nil {
		t.Fatalf("Malloc failed: %v", err)
	}
	// 验证 Malloc 后 Redis 是否有值
	// 注意：Malloc 返回的是分配的起始值，Redis 中存储的是 CURR (当前分配到的值)
	// 第一次 Malloc(10): DB初始化0->60(预分配50+10), Redis CURR=10, LAST=60

	// 场景 2: 缓存中不存在 (DB 中有数据，但 Redis 被清空)
	// 先给 conv_002 在 DB 中造数据
	db.Create(&model.SeqConversation{ID: convID2, MaxSeq: 100})
	// 此时 Redis 中没有 conv_002 的 key

	// 场景 3: 缓存和 DB 都不存在 (conv_003)
	// 预期应该返回 0 或初始化为 0

	// 执行批量获取
	convIDs := []string{convID1, convID2, convID3}
	seqs, err := seqCache.GetMaxSeqs(ctx, convIDs)
	if err != nil {
		t.Fatalf("GetMaxSeqs failed: %v", err)
	}

	// 验证结果
	// 1. conv_001 (缓存存在): 应该是 10 (Malloc 申请了 10)
	if seq, ok := seqs[convID1]; !ok || seq != 10 {
		t.Errorf("convID1 seq expected 10, got %d (exists: %v)", seq, ok)
	}

	// 2. conv_002 (缓存不存在，DB存在):
	// GetMaxSeqs -> getSingleMaxSeq -> getMaxSeq -> Malloc(size=0)
	// Malloc(0) 会从 DB 加载。DB 中是 100。
	// 逻辑：Malloc(0) -> Redis不存在 -> MallocFromDB(size=0) -> DB=100 -> Redis CURR=100
	if seq, ok := seqs[convID2]; !ok || seq != 100 {
		t.Errorf("convID2 seq expected 100, got %d (exists: %v)", seq, ok)
	}

	// 3. conv_003 (都不存在):
	// Malloc(0) -> Redis不存在 -> MallocFromDB(size=0) -> DB Insert 0 -> Redis CURR=0
	if seq, ok := seqs[convID3]; !ok || seq != 0 {
		t.Errorf("convID3 seq expected 0, got %d (exists: %v)", seq, ok)
	}

	t.Logf("GetMaxSeqs result: %+v", seqs)
}
