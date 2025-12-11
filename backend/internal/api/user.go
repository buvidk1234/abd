package api

import (
	"backend/internal/api/apiresp"
	"backend/internal/api/apiresp/errs"
	"backend/internal/service"
	"backend/pkg/util"
	"strconv"

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
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	userID, err := a.userService.UserRegister(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, strconv.FormatInt(userID, 10))
}

// UpdateUserInfo 更新用户信息
func (a *UserApi) UpdateUserInfo(c *gin.Context) {
	var req service.UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	err := a.userService.UpdateUserInfo(c.Request.Context(), req, c.GetInt64("user_id"))
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

func (a *UserApi) GetUsersPublicInfo(c *gin.Context) {
	userInfo, err := a.userService.GetUsersPublicInfo(c.Request.Context(), c.GetInt64("user_id"))
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, userInfo)
}

func (a *UserApi) SearchUser(c *gin.Context) {
	var req service.SearchUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	users, err := a.userService.SearchUser(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, users)
}

func (a *UserApi) UserLogin(c *gin.Context) {
	var req service.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	token, err := a.userService.UserLogin(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, token)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization Header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort() // 终止请求
			return
		}

		// 解析 Token
		user_id, err := util.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort() // 终止请求
			return
		}

		// 设置用户 ID 到上下文
		c.Set("user_id", user_id)
		c.Next() // 继续处理请求
	}
}
