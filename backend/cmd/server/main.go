package main

import (
	"backend/internal/api"
	"backend/internal/model"
	"backend/internal/pkg/database"
)

func main() {
	database.Init()
	db := database.GetDB()
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Friend{})
	db.AutoMigrate(&model.FriendRequest{})

	r := api.NewGinRouter()

	r.Run()
}
