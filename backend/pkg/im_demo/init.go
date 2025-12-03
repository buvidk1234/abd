package im_demo

import "github.com/gin-gonic/gin"

func Init(gin *gin.Engine) {
	gin.GET("/ws", WsHandler)
}
