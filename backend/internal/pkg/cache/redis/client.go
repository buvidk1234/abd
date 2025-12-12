package redis

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var (
	RDB *redis.Client
	sf  singleflight.Group
)

func Init(cfg Config) *redis.Client {
	RDB = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  3 * time.Second,
	})

	// 启动时 ping，失败直接报错
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		log.Printf("redis: failed to connect to redis at %s: %v", cfg.Addr, err)
		panic(err)
	}
	return RDB
}

func GetRDB() *redis.Client {
	if RDB == nil {
		log.Panic("redis: RDB is not initialized, call Init first")
	}
	return RDB
}

func GetCacheString(key string, fn func() (string, error), expire time.Duration) (string, error) {
	ctx := context.Background()
	val, err := RDB.Get(ctx, key).Result()
	if err == nil {
		return val, nil
	}
	v, err, _ := sf.Do(key, func() (interface{}, error) {
		res, err := fn()
		if err != nil {
			return "", err
		}
		var ex time.Duration
		if res == "" {
			ex = 1 * time.Minute
		} else {
			ex = expire - time.Duration(rand.Float64()*0.1*float64(expire))
		}

		if err := RDB.Set(ctx, key, res, ex).Err(); err != nil {
			log.Printf("redis set failed: %v", err)
		}
		return res, nil
	})
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func GetCache[T any](key string, fn func() (T, error), expire time.Duration) (T, error) {
	var zero T
	ctx := context.Background()
	val, err := RDB.Get(ctx, key).Result()
	if err == nil {
		// 如果 T 是 string，直接返回，避免不必要的 JSON Unmarshal
		if v, ok := any(val).(T); ok {
			return v, nil
		}
		// 否则尝试反序列化
		var res T
		err = json.Unmarshal([]byte(val), &res)
		if err == nil {
			return res, nil
		}
		// 反序列化失败，视为缓存未命中，继续向下执行
	}

	v, err, _ := sf.Do(key, func() (interface{}, error) {
		res, err := fn()
		if err != nil {
			return zero, err
		}

		var valStr string
		var isEmpty bool

		// 序列化逻辑
		if s, ok := any(res).(string); ok {
			valStr = s
			isEmpty = (s == "")
		} else {
			b, err := json.Marshal(res)
			if err != nil {
				return zero, err
			}
			valStr = string(b)
			isEmpty = (valStr == "null")
		}

		var ex time.Duration
		if isEmpty {
			ex = 1 * time.Minute
		} else {
			ex = expire - time.Duration(rand.Float64()*0.1*float64(expire))
		}

		if err := RDB.Set(ctx, key, valStr, ex).Err(); err != nil {
			log.Printf("redis set failed: %v", err)
		}
		return res, nil
	})
	if err != nil {
		return zero, err
	}
	return v.(T), nil
}

func BatchGetCache[T any](keys []string, fn func([]string) (map[string]T, error), expire time.Duration) (map[string]T, error) {
	ctx := context.Background()
	result := make(map[string]T)

	if len(keys) == 0 {
		return result, nil
	}

	vals, err := RDB.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var missingKeys []string
	for i, v := range vals {
		if v == nil {
			missingKeys = append(missingKeys, keys[i])
		} else {
			if s, ok := v.(string); ok {
				// Try to unmarshal if T is not string
				var t T
				if _, isString := any(t).(string); isString {
					result[keys[i]] = any(s).(T)
				} else {
					if err := json.Unmarshal([]byte(s), &t); err == nil {
						result[keys[i]] = t
					} else {
						// Unmarshal failed, treat as missing
						missingKeys = append(missingKeys, keys[i])
					}
				}
			} else {
				missingKeys = append(missingKeys, keys[i])
			}
		}
	}

	if len(missingKeys) > 0 {
		dbRes, err := fn(missingKeys)
		if err != nil {
			return nil, err
		}

		pipe := RDB.Pipeline()
		for _, key := range missingKeys {
			val := dbRes[key]
			result[key] = val

			var valStr string
			var isEmpty bool

			if s, ok := any(val).(string); ok {
				valStr = s
				isEmpty = (s == "")
			} else {
				b, err := json.Marshal(val)
				if err != nil {
					// Skip this key if marshal fails, or handle error
					continue
				}
				valStr = string(b)
				isEmpty = (valStr == "null")
			}

			var ex time.Duration
			if isEmpty {
				ex = 1 * time.Minute
			} else {
				ex = expire - time.Duration(rand.Float64()*0.1*float64(expire))
			}
			pipe.Set(ctx, key, valStr, ex)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			log.Printf("redis batch set failed: %v", err)
		}
	}

	return result, nil
}
