package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SecretPassword = "unimate_secret_password_2024"

func GenToken(id uint) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	jwt_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token, _ := jwt_token.SignedString([]byte(SecretPassword))
	return token
}
