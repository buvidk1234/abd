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

	r := api.NewGinRouter()

	r.Run()
}
