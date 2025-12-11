package redis

const (
	UserTokenKey = "auth:user:%s:%s" // userID, platform
	TokenUserKey = "auth:token:%s"   // token
)
