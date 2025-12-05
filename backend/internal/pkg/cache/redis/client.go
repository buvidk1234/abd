package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func Init(cfg Config) error {
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
	if err := RDB.Ping(Ctx).Err(); err != nil {
		return err
	}
	return nil
}

func GetRDB() *redis.Client {
	return RDB
}
