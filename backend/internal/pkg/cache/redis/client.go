package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
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
