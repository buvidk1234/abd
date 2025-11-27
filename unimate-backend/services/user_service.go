package services

import (
	"unimate-backend/database"
	"unimate-backend/models"
	"unimate-backend/utils"
)

func Register(username, email, password string) error {
	db := database.GetDB()

	user := models.User{
		Username: username,
		Email:    email,
		Password: utils.HashPassword(password),
	}

	return db.Create(&user).Error
}

func Login(email, password string) (string, error) {
	db := database.GetDB()

	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(user.Password, password) {
		return "", err
	}

	return utils.GenerateToken(user.ID)
}
