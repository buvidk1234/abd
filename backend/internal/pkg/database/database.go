package database

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// Init 初始化数据库（只执行一次）
func Init() *gorm.DB {
	once.Do(func() {
		var err error

		DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		// 连接池设置
		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatalf("failed to get generic database object: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		log.Println("database initialized successfully")
	})

	return DB
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		return Init()
	}
	return DB
}
