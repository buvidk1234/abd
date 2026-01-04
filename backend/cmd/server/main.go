package main

import (
	"backend/internal/api"
	"backend/internal/im"
	"backend/internal/im/distributor"
	"backend/internal/im/pusher"
	"backend/internal/model"
	"backend/internal/pkg/cache/redis"
	"backend/internal/pkg/database"
	"backend/internal/pkg/kafka"
	"backend/internal/pkg/prommetrics"
	"backend/internal/pkg/snowflake"
	"context"
	"os"

	"github.com/goccy/go-yaml"
)

type AppConfig struct {
	Redis     redis.Config     `yaml:"redis"`
	Snowflake snowflake.Config `yaml:"snowflake"`
	Kafka     kafka.Config     `yaml:"kafka"`
	Database  database.Config  `yaml:"database"`
	Server    ServerConfig     `yaml:"server"`
	WebSocket im.Config        `yaml:"websocket"`
}

type ServerConfig struct {
	HTTPAddr    string `yaml:"http_addr"`
	MetricsAddr string `yaml:"metrics_addr"`
}

func main() {

	var cfg AppConfig
	raw, _ := os.ReadFile("config/config.yaml")
	yaml.Unmarshal(raw, &cfg)

	redis.Init(cfg.Redis)
	database.Init(cfg.Database)
	snowflake.Init(cfg.Snowflake)
	kafka.Init(cfg.Kafka)

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
	wsServer := im.NewWsServer(cfg.WebSocket)
	pusher.InitAndRun(wsServer)
	distributor := distributor.NewDistributor(wsServer)
	go distributor.Start()
	go wsServer.Run(context.Background())

	prommetrics.RegistryAll()
	go prommetrics.Start(cfg.Server.MetricsAddr)

	r.Run(cfg.Server.HTTPAddr)
}
