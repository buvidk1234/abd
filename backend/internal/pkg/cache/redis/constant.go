package redis

import "time"

const (
	UserTokenKey = "auth:user:%s:%s"    // userID, platform
	TokenUserKey = "auth:token:%s"      // token
	ExpireTime   = 365 * 24 * time.Hour // TODO: 不同业务可能需要不同的过期时间，后续可调整
)
