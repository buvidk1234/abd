package utils

import (
	"time"
	"unimate-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(config.JWTSecret)
}
