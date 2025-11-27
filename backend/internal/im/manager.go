package im

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// 1. 定义全局管理器
var IMManager = &Manager{
	Clients: make(map[string]*Client),
}

type Manager struct {
	Clients    map[string]*Client // 用户ID -> 连接对象 的映射
	Lock       sync.RWMutex       // 读写锁，保证并发安全
	Register   chan *Client       // 注册通道 (可选优化，这里直接用方法调用)
	Unregister chan *Client       // 注销通道
}

// 2. 注册连接 (用户上线)
func (m *Manager) AddClient(uid string, conn *websocket.Conn) {
	client := &Client{
		ID:   uid,
		Conn: conn,
		Send: make(chan []byte, 256), // 带缓冲的发送管道，防止阻塞
	}

	m.Lock.Lock()
	// 如果该用户之前有旧连接，先关掉旧的 (踢下线逻辑)
	if oldClient, ok := m.Clients[uid]; ok {
		close(oldClient.Send)
	}
	m.Clients[uid] = client
	m.Lock.Unlock()

	log.Printf("[IM] 用户 %s 上线", uid)

	// 启动协程处理读写
	go client.WritePump()
	go client.ReadPump()
}

// 3. 注销连接 (用户下线/断开)
func (m *Manager) RemoveClient(uid string) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if client, ok := m.Clients[uid]; ok {
		close(client.Send)     // 关闭通道，通知 WritePump 退出
		client.Conn.Close()    // 关闭 Socket 连接
		delete(m.Clients, uid) // 从 Map 中移除
		log.Printf("[IM] 用户 %s 下线", uid)
	}
}

// 4. 发送消息给指定用户
func (m *Manager) SendToUser(toUserID string, msg []byte) bool {
	m.Lock.RLock()
	targetClient, ok := m.Clients[toUserID]
	m.Lock.RUnlock()

	if ok {
		// 用户在线，放入他的发送管道
		select {
		case targetClient.Send <- msg:
			return true
		default:
			// 管道满了 (比如对方网太卡)，这里可以选择丢弃或者强行断开
			return false
		}
	}

	// 用户不在线，返回 false
	// 调用者应该负责将消息存入数据库，标记为未读
	return false
}

// 5. 简单的消息分发逻辑 (业务层)
func Dispatch(packet *MsgPacket) {
	// --- TODO: 在这里插入保存到 MySQL 的代码 ---
	// DB.Save(packet)
	// ---------------------------------------

	log.Printf("[IM] 转发消息: %s -> %s: %s", packet.FromID, packet.ToID, packet.Content)

	// 转成 JSON 字节流
	msgBytes, _ := json.Marshal(packet)

	// 尝试发送给接收者
	isSent := IMManager.SendToUser(packet.ToID, msgBytes)

	if !isSent {
		log.Printf("[IM] 用户 %s 不在线，消息已存库(模拟)", packet.ToID)
	}
}
