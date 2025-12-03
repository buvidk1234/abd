package im_demo

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string          // 用户ID
	Conn *websocket.Conn // WebSocket 连接
	Send chan []byte     // 待发送的数据管道
}

// WritePump: 专门负责【写】 (后端 -> 前端)
// 这是一个死循环，监听 Send 管道，有数据就发给前端
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		// 批量收集所有待发送消息，减少锁和系统调用
		messages := [][]byte{message}
		for {
			select {
			case m, ok := <-c.Send:
				if !ok {
					break
				}
				messages = append(messages, m)
			default:
				goto SEND
			}
		}
	SEND:
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		for _, msg := range messages {
			if _, err := w.Write(msg); err != nil {
				w.Close()
				return
			}
		}
		if err := w.Close(); err != nil {
			return
		}
	}
	// 管道被关闭 (RemoveClient 被调用)
	_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// ReadPump: 专门负责【读】 (前端 -> 后端)
// 这是一个死循环，阻塞读取前端发来的 JSON
func (c *Client) ReadPump() {
	defer func() {
		IMManager.RemoveClient(c.ID) // 出错或退出时，注销自己
		c.Conn.Close()
	}()

	// 设置读取限制 (防止恶意发送超大包)
	c.Conn.SetReadLimit(4096)
	// 设置心跳超时 (如果60秒没收到消息，视为断开)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// 设置收到 Pong (心跳回应) 的处理
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var packet MsgPacket
		// 读取 JSON
		err := c.Conn.ReadJSON(&packet)
		if err != nil {
			// 这种错误通常是连接断开 (websocket: close 1006)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 补全信息 (强制修正 FromID，防止伪造)
		packet.FromID = c.ID
		packet.Seq = time.Now().UnixMilli()

		// 处理心跳包 (Type=0)
		if packet.Type == MsgTypeHeartbeat {
			// 收到前端的心跳，刷新超时时间，不进行分发
			c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			// 可选：回复一个 Pong
			// c.Send <- []byte(`{"type":0}`)
			continue
		}

		// 将消息丢给 Manager 的 Dispatch 进行分发或存库
		Dispatch(&packet)
	}
}
