package main

import (
	"backend/internal/api"
	"backend/internal/im"
	"backend/internal/model"
	"backend/internal/pkg/cache/redis"
	"backend/internal/pkg/database"
	"backend/internal/pkg/snowflake"
	"context"
	"os"

	"github.com/goccy/go-yaml"
)

type AppConfig struct {
	Redis     redis.Config     `yaml:"redis"`
	Snowflake snowflake.Config `yaml:"snowflake"`
}

func main() {

	var cfg AppConfig
	raw, _ := os.ReadFile("config/config.yaml")
	yaml.Unmarshal(raw, &cfg)

	redis.Init(cfg.Redis)
	database.Init()
	snowflake.Init(cfg.Snowflake)

	db := database.GetDB()
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Friend{})
	db.AutoMigrate(&model.FriendRequest{})
	db.AutoMigrate(&model.Group{})
	db.AutoMigrate(&model.GroupRequest{})
	db.AutoMigrate(&model.GroupMember{})

	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Conversation{})
	db.AutoMigrate(&model.SeqConversation{})
	db.AutoMigrate(&model.SeqUser{})
	db.AutoMigrate(&model.UserTimeline{})

	r := api.NewGinRouter()
	wsServer := im.NewWsServer()
	go wsServer.Run(context.Background())

	r.Run()
}
