package api

import (
	"backend/internal/service"

	"backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ws", WsHandler)

	// middlewares
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// user routes
	u := NewUserApi(service.NewUserService(database.GetDB()))
	{
		userRouterGroup := r.Group("/user")
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", u.UpdateUserInfo)
		userRouterGroup.POST("/get_users_info", u.GetUsersPublicInfo)
		userRouterGroup.POST("/user_login", u.UserLogin)
	}

	return r
}
