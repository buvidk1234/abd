package redis

import (
	"backend/internal/model"
	"backend/internal/pkg/database"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
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
	seqGenerator := NewSeqConversationCacheRedis(database.GetDB(), RDB)
	RDB.FlushDB(context.Background())
	db := database.GetDB()
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
