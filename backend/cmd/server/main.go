package main

import (
	"backend/internal/api"
	"backend/internal/model"
	"backend/internal/pkg/database"
	"backend/internal/pkg/snowflake"
)

func main() {
	database.Init()
	snowflake.Init(1)

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

	r.Run()
}
