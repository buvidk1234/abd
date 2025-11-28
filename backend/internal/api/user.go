package api

import (
	"backend/internal/service"

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
	userDTOs, err := a.userService.GetUsersPublicInfo(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"users": userDTOs})
}

func (a *UserApi) UserLogin(c *gin.Context) {
	var req service.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	err := a.userService.UserLogin(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "登录成功"})
}
