package api

import (
	"backend/internal/service"
	"backend/pkg/util"

	"github.com/gin-gonic/gin"
)

type UserApi struct {
	userService *service.UserService
}

func NewUserApi(userService *service.UserService) *UserApi {
	return &UserApi{userService: userService}
}

// UserRegister 用户注册
func (a *UserApi) UserRegister(c *gin.Context) {
	var req service.UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	userID, err := a.userService.UserRegister(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"userID": userID})
}

// UpdateUserInfo 更新用户信息
func (a *UserApi) UpdateUserInfo(c *gin.Context) {
	var req service.UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	err := a.userService.UpdateUserInfo(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "更新成功"})
}

func (a *UserApi) GetUsersPublicInfo(c *gin.Context) {
	var req service.GetUsersPublicInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	userInfos, err := a.userService.GetUsersPublicInfo(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"userInfos": userInfos})
}

func (a *UserApi) UserLogin(c *gin.Context) {
	var req service.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	token, err := a.userService.UserLogin(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 放行登录路径
		if c.Request.URL.Path == "/user/user_login" {
			c.Next() // 直接跳过中间件逻辑
			return
		}

		// 获取 Authorization Header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization token is required"})
			c.Abort() // 终止请求
			return
		}

		// 解析 Token
		my_user_id, err := util.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort() // 终止请求
			return
		}

		// 设置用户 ID 到上下文
		c.Set("my_user_id", my_user_id)
		c.Next() // 继续处理请求
	}
}
