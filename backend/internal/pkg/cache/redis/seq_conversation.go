package redis

import (
	"backend/internal/pkg/cache/cachekey"
	"context"
	"fmt"
	"strings"
	"time"

	"backend/internal/model"

	"log"

	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeqConversationCacheRedis struct {
	db               *gorm.DB
	client           *redis.Client
	lockTime         time.Duration
	dataTime         time.Duration
	minSeqExpireTime time.Duration
}

func NewSeqConversationCacheRedis(db *gorm.DB, client *redis.Client) *SeqConversationCacheRedis {
	return &SeqConversationCacheRedis{
		db:               db,
		lockTime:         time.Second * 3,
		dataTime:         time.Hour * 24 * 365,
		minSeqExpireTime: time.Hour,
		client:           client,
	}
}

func (s *SeqConversationCacheRedis) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	minSeq, err := GetCache(cachekey.GetSeqConvMinSeqKey(conversationID), func() (int64, error) {
		var seqConv model.SeqConversation
		if err := s.db.WithContext(ctx).Where("id = ?", conversationID).First(&seqConv).Error; err != nil {
			return 0, err
		}
		return seqConv.MinSeq, nil
	}, s.minSeqExpireTime)
	return minSeq, err
}

func (s *SeqConversationCacheRedis) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	switch len(conversationIDs) {
	case 0:
		return map[string]int64{}, nil
	case 1:
		return s.getSingleMaxSeq(ctx, conversationIDs[0])
	}
	keys := make([]string, 0, len(conversationIDs))
	keyConversationID := make(map[string]string, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		key := cachekey.GetSeqConvKey(conversationID)
		if _, ok := keyConversationID[key]; ok {
			continue
		}
		keys = append(keys, key)
		keyConversationID[key] = conversationID
	}
	if len(keys) == 1 {
		return s.getSingleMaxSeq(ctx, conversationIDs[0])
	}

	seqs := make(map[string]int64, len(conversationIDs))
	if err := s.batchGetMaxSeq(ctx, keys, keyConversationID, seqs); err != nil {
		return nil, err
	}

	return seqs, nil
}

func (s *SeqConversationCacheRedis) batchGetMaxSeq(ctx context.Context, keys []string, keyConversationID map[string]string, seqs map[string]int64) error {
	result := make([]*redis.StringCmd, len(keys))
	pipe := s.client.Pipeline()
	for i, key := range keys {
		result[i] = pipe.HGet(ctx, key, "CURR")
	}
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	var notFoundKey []string
	for i, r := range result {
		req, err := r.Int64()
		if err == nil {
			seqs[keyConversationID[keys[i]]] = req
		} else if errors.Is(err, redis.Nil) {
			notFoundKey = append(notFoundKey, keys[i])
		} else {
			return err
		}
	}
	for _, key := range notFoundKey {
		conversationID := keyConversationID[key]
		seq, err := s.GetMaxSeq(ctx, conversationID)
		if err != nil {
			return err
		}
		seqs[conversationID] = seq
	}
	return nil
}

func (s *SeqConversationCacheRedis) getSingleMaxSeq(ctx context.Context, conversationID string) (map[string]int64, error) {
	seq, err := s.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return map[string]int64{conversationID: seq}, nil
}

// 获取当前最大序列号
func (s *SeqConversationCacheRedis) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return s.Malloc(ctx, conversationID, 0)
}

func (s *SeqConversationCacheRedis) Malloc(ctx context.Context, conversationID string, size int64) (int64, error) {
	seq, _, err := s.mallocTime(ctx, conversationID, size)
	return seq, err
}

func (s *SeqConversationCacheRedis) mallocTime(ctx context.Context, conversationID string, size int64) (int64, int64, error) {
	if size < 0 {
		return 0, 0, errors.New("size must be greater than 0")
	}
	key := cachekey.GetSeqConvKey(conversationID)
	for i := 0; i < 10; i++ {
		states, err := s.malloc(ctx, key, size)
		if err != nil {
			return 0, 0, err
		}
		switch states[0] {
		case 0: // success
			return states[1], states[3], nil
		case 1: // not found
			mallocSize := s.getMallocSize(conversationID, size)
			seq, err := s.mallocFromDB(ctx, conversationID, mallocSize) // (0,5]
			if err != nil {
				return 0, 0, err
			}
			s.setSeqRetry(ctx, key, states[1], seq+size, seq+mallocSize, states[2])
			return seq, 0, nil
		case 2: // locked
			if err := s.wait(ctx); err != nil {
				return 0, 0, err
			}
			continue
		case 3: // exceeded cache max value
			currSeq := states[1]
			lastSeq := states[2]
			mill := states[4]
			mallocSize := s.getMallocSize(conversationID, size)
			seq, err := s.mallocFromDB(ctx, conversationID, mallocSize)
			if err != nil {
				return 0, 0, err
			}
			if lastSeq == seq {
				s.setSeqRetry(ctx, key, states[3], currSeq+size, seq+mallocSize, mill)
				return currSeq, states[4], nil
			} else {
				log.Printf("malloc seq not equal cache last seq conversationID=%s currSeq=%d lastSeq=%d mallocSeq=%d", conversationID, currSeq, lastSeq, seq)
				s.setSeqRetry(ctx, key, states[3], seq+size, seq+mallocSize, mill)
				return seq, mill, nil
			}
		default:
			log.Printf("malloc seq unknown state state=%d conversationID=%s size=%d", states[0], conversationID, size)
			return 0, 0, fmt.Errorf("unknown state: %d", states[0])
		}
	}
	log.Printf("malloc seq retrying still failed conversationID=%s size=%d", conversationID, size)
	return 0, 0, fmt.Errorf("malloc seq waiting for lock timeout conversationID=%s size=%d", conversationID, size)
}

func (s *SeqConversationCacheRedis) mallocFromDB(ctx context.Context, conversationID string, size int64) (int64, error) {
	// newseq is distributed max seq, (max,+inf)
	var seqModel model.SeqConversation
	var oldSeq int64
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).FirstOrCreate(&seqModel, model.SeqConversation{ID: conversationID, MaxSeq: 0}).Error; err != nil {
			return err
		}
		oldSeq = seqModel.MaxSeq
		seqModel.MaxSeq = seqModel.MaxSeq + size
		if err := tx.Model(&seqModel).Where("id = ?", conversationID).Update("max_seq", seqModel.MaxSeq).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return oldSeq, nil
}

// malloc size=0 is to get the current seq size>0 is to allocate seq
func (s *SeqConversationCacheRedis) malloc(ctx context.Context, key string, size int64) ([]int64, error) {
	/*
		1. key不存在，创建，加锁，去申请序列号
		2. key存在且被加锁，有线程在申请序列号，等待
		3. key存在但序列号不够用，加锁，去申请序列号
		4. key存在且序列号够用，分配，返回
		key: {
			"CURR": 当前分配到的序列号
			"LAST": 最大序列号
			"LOCK": 锁标识
			"TIME": 分配时间戳
		}
	*/
	// 0: success
	// 1: need to obtain and lock
	// 2: already locked
	// 3: exceeded the maximum value and locked
	script := `
local key = KEYS[1]
local size = tonumber(ARGV[1])
local lockSecond = ARGV[2]
local dataSecond = ARGV[3]
local mallocTime = ARGV[4]
local result = {}
if redis.call("EXISTS", key) == 0 then
	local lockValue = math.random(0, 999999999)
	redis.call("HSET", key, "LOCK", lockValue)
	redis.call("EXPIRE", key, lockSecond)
	table.insert(result, 1)
	table.insert(result, lockValue)
	table.insert(result, mallocTime)
	return result
end
if redis.call("HEXISTS", key, "LOCK") == 1 then
	table.insert(result, 2)
	return result
end
local curr_seq = tonumber(redis.call("HGET", key, "CURR"))
local last_seq = tonumber(redis.call("HGET", key, "LAST"))
if size == 0 then
	redis.call("EXPIRE", key, dataSecond)
	table.insert(result, 0)
	table.insert(result, curr_seq)
	table.insert(result, last_seq)
	local setTime = redis.call("HGET", key, "TIME")
	if setTime then
		table.insert(result, setTime)	
	else
		table.insert(result, 0)
	end
	return result
end
local max_seq = curr_seq + size
if max_seq > last_seq then
	local lockValue = math.random(0, 999999999)
	redis.call("HSET", key, "LOCK", lockValue)
	redis.call("HSET", key, "CURR", last_seq)
	redis.call("HSET", key, "TIME", mallocTime)
	redis.call("EXPIRE", key, lockSecond)
	table.insert(result, 3)
	table.insert(result, curr_seq)
	table.insert(result, last_seq)
	table.insert(result, lockValue)
	table.insert(result, mallocTime)
	return result
end
redis.call("HSET", key, "CURR", max_seq)
redis.call("HSET", key, "TIME", ARGV[4])
redis.call("EXPIRE", key, dataSecond)
table.insert(result, 0)
table.insert(result, curr_seq)
table.insert(result, last_seq)
table.insert(result, mallocTime)
return result
`
	result, err := s.client.Eval(ctx, script, []string{key}, size, int64(s.lockTime/time.Second), int64(s.dataTime/time.Second), time.Now().UnixMilli()).Int64Slice()
	if err != nil {
		return nil, fmt.Errorf("redis eval failed: %w", err)
	}
	return result, nil
}

func (s *SeqConversationCacheRedis) getMallocSize(conversationID string, size int64) int64 {
	if size == 0 {
		return 0
	}
	var basicSize int64
	if IsGroupConversationID(conversationID) {
		basicSize = 100
	} else {
		basicSize = 50
	}
	basicSize += size
	return basicSize
}

func IsGroupConversationID(conversationID string) bool {
	return strings.HasPrefix(conversationID, "g")
}

func (s *SeqConversationCacheRedis) setSeqRetry(ctx context.Context, key string, owner int64, currSeq int64, lastSeq int64, mill int64) {
	for i := 0; i < 10; i++ {
		state, err := s.setSeq(ctx, key, owner, currSeq, lastSeq, mill)
		if err != nil {
			log.Printf("set seq cache failed key=%s owner=%d currSeq=%d lastSeq=%d attempt=%d err=%v", key, owner, currSeq, lastSeq, i+1, err)
			if err := s.wait(ctx); err != nil {
				return
			}
			continue
		}
		switch state {
		case 0: // ideal state
		case 1:
			log.Printf("set seq cache lock not found key=%s owner=%d currSeq=%d lastSeq=%d", key, owner, currSeq, lastSeq)
		case 2:
			log.Printf("set seq cache lock held by someone else key=%s owner=%d currSeq=%d lastSeq=%d", key, owner, currSeq, lastSeq)
		default:
			log.Printf("set seq cache lock unknown state key=%s owner=%d currSeq=%d lastSeq=%d", key, owner, currSeq, lastSeq)
		}
		return
	}
	log.Printf("set seq cache retrying still failed key=%s owner=%d currSeq=%d lastSeq=%d", key, owner, currSeq, lastSeq)
}

func (s *SeqConversationCacheRedis) wait(ctx context.Context) error {
	timer := time.NewTimer(time.Second / 4)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *SeqConversationCacheRedis) setSeq(ctx context.Context, key string, owner int64, currSeq int64, lastSeq int64, mill int64) (int64, error) {
	if lastSeq < currSeq {
		return 0, errors.New("lastSeq must be greater than currSeq")
	}
	// 0: success
	// 1: success the lock has expired, but has not been locked by anyone else
	// 2: already locked, but not by yourself
	script := `
local key = KEYS[1]
local lockValue = ARGV[1]
local dataSecond = ARGV[2]
local curr_seq = tonumber(ARGV[3])
local last_seq = tonumber(ARGV[4])
local mallocTime = ARGV[5]
if redis.call("EXISTS", key) == 0 then
	redis.call("HSET", key, "CURR", curr_seq, "LAST", last_seq, "TIME", mallocTime)
	redis.call("EXPIRE", key, dataSecond)
	return 1
end
if redis.call("HGET", key, "LOCK") ~= lockValue then
	return 2
end
redis.call("HDEL", key, "LOCK")
redis.call("HSET", key, "CURR", curr_seq, "LAST", last_seq, "TIME", mallocTime)
redis.call("EXPIRE", key, dataSecond)
return 0
`
	result, err := s.client.Eval(ctx, script, []string{key}, owner, int64(s.dataTime/time.Second), currSeq, lastSeq, mill).Int64()
	if err != nil {
		return 0, fmt.Errorf("redis eval failed: %w", err)
	}
	return result, nil
}
