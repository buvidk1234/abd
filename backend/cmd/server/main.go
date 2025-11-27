package main

import (
	"backend/internal/api" // 你的 api 包路径

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// ... 其他 HTTP 路由 (如登录注册) ...

	// 【新增】WebSocket 路由
	// 注意：WS 通常用 GET 请求，参数在 URL 里 (?token=xxx)
	r.GET("/ws", api.WsHandler)

	return r
}

func main() {
	r := InitRouter()
	r.Run()
}
