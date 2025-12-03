package im_demo

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理 WebSocket 请求
// 路由: GET /ws?token=xxxxx
func WsHandler(c *gin.Context) {
	// 1. 获取 Token 并解析出 UserID (鉴权)
	token := c.Query("token")
	// userID, err := utils.ParseToken(token) // 你的 JWT 解析逻辑
	// if err != nil {
	// 	c.JSON(401, gin.H{"msg": "Unauthorized"})
	// 	return
	// }
	userID := token // 临时简化处理，直接用 token 作为 userID

	// 2. 升级 HTTP 为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// 3. 将连接交给 IM Manager 管理
	// 这一步之后，ReadPump 和 WritePump 就会自动开始工作
	IMManager.AddClient(userID, conn)
}
