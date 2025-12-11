package redis

import (
	"context"
	"fmt"
	"time"
)

type TokenRepository struct{}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (r *TokenRepository) SetUserToken(ctx context.Context, userID, platform, token string, ttl time.Duration) error {
	key := fmt.Sprintf(UserTokenKey, userID, platform)
	return RDB.Set(ctx, key, token, ttl).Err()
}

func (r *TokenRepository) SetTokenUser(ctx context.Context, token, userID, platform string, ttl time.Duration) error {
	key := fmt.Sprintf(TokenUserKey, token)
	val := fmt.Sprintf("%s:%s", userID, platform)
	return RDB.Set(ctx, key, val, ttl).Err()
}

func (r *TokenRepository) GetTokenByUser(ctx context.Context, userID, platform string) (string, error) {
	key := fmt.Sprintf(UserTokenKey, userID, platform)
	return RDB.Get(ctx, key).Result()
}

func (r *TokenRepository) GetUserByToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf(TokenUserKey, token)
	return RDB.Get(ctx, key).Result()
}

func (r *TokenRepository) DeleteUserToken(ctx context.Context, userID, platform string) error {
	key := fmt.Sprintf(UserTokenKey, userID, platform)
	return RDB.Del(ctx, key).Err()
}

func (r *TokenRepository) DeleteTokenUser(ctx context.Context, token string) error {
	key := fmt.Sprintf(TokenUserKey, token)
	return RDB.Del(ctx, key).Err()
}
