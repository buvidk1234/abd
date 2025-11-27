package middlewares

import (
	"net/http"
	"unimate-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !t.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
