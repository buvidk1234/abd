package database

import (
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// Config 数据库配置
type Config struct {
	Driver       string `yaml:"driver"`         // sqlite3 | mysql
	DSN          string `yaml:"dsn"`            // 数据源连接字符串
	MaxOpenConns int    `yaml:"max_open_conns"` // 最大打开连接数
	MaxIdleConns int    `yaml:"max_idle_conns"` // 最大空闲连接数
}

// Init 初始化数据库（只执行一次）
func Init(cfg Config) *gorm.DB {
	once.Do(func() {
		var err error

		switch cfg.Driver {
		case "mysql":
			DB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
		case "sqlite3", "sqlite":
			dsn := cfg.DSN
			if dsn == "" {
				dsn = "gorm.db"
			}
			DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		default:
			log.Fatalf("unsupported database driver: %s", cfg.Driver)
		}

		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		// 连接池设置
		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatalf("failed to get generic database object: %v", err)
		}

		if cfg.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		} else {
			sqlDB.SetMaxIdleConns(10)
		}

		if cfg.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		} else {
			sqlDB.SetMaxOpenConns(100)
		}

		log.Println("database initialized successfully")
	})

	return DB
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("database not initialized, call Init() first")
	}
	return DB
}
