package redis

import (
	"backend/internal/model"
	"backend/internal/pkg/cache/cachekey"
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SeqUserCacheRedis struct {
	db     *gorm.DB
	client *redis.Client
}

func NewSeqUserCacheRedis(db *gorm.DB, client *redis.Client) *SeqUserCacheRedis {
	return &SeqUserCacheRedis{
		db:     db,
		client: client,
	}
}

func (s *SeqUserCacheRedis) GetSeqUserMinSeq(ctx context.Context, userID int64, conversationID string) (int64, error) {
	key := cachekey.GetSeqUserMinSeqKey(strconv.FormatInt(userID, 10), conversationID)
	userMinSeq, err := GetCache(key, func() (int64, error) {
		var seqUser model.SeqUser
		if err := s.db.WithContext(ctx).Find(&seqUser, "user_id = ? AND conversation_id = ?", userID, conversationID).Error; err != nil {
			return 0, err
		}
		return seqUser.MinSeq, nil
	}, ExpireTime)
	if err != nil {
		return 0, err
	}
	return userMinSeq, nil
}

func (s *SeqUserCacheRedis) GetSeqUserMaxSeq(ctx context.Context, userID int64, conversationID string) (int64, error) {
	key := cachekey.GetSeqUserMaxSeqKey(strconv.FormatInt(userID, 10), conversationID)
	userMaxSeq, err := GetCache(key, func() (int64, error) {
		var seqUser model.SeqUser
		if err := s.db.WithContext(ctx).Find(&seqUser, "user_id = ? AND conversation_id = ?", userID, conversationID).Error; err != nil {
			return 0, err
		}
		return seqUser.MaxSeq, nil
	}, ExpireTime)
	if err != nil {
		return 0, err
	}
	return userMaxSeq, nil
}
