package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "my_secret_key"

// 定义内部使用的 Claims 结构体
type userClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 Token
// 默认过期时间设置为 24 小时
func GenerateToken(userID string) (string, error) {
	claims := userClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间：当前时间 + 24小时
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// 设置签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 设置生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用 HS256 签名算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 进行签名
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析 Token 并返回 UserID
func ParseToken(tokenStr string) (string, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	// 验证并提取数据
	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token")
}
