package main

import (
	"unimate-backend/database"
	"unimate-backend/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.Init()

	r := gin.Default()

	// 注册路由
	routers.SetupRouter(r)

	// 启动服务
	r.Run(":8080")
}
