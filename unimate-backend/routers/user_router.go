package routers

import (
	"unimate-backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	user := r.Group("/users")
	{
		user.POST("/register", handlers.Register)
		user.POST("/login", handlers.Login)
	}
}
