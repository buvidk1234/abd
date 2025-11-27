package im

// 消息类型常量
const (
	MsgTypeHeartbeat = 0 // 心跳包
	MsgTypeText      = 1 // 文本消息
	MsgTypeImage     = 2 // 图片消息 (Content为URL)
)

// MsgPacket 核心消息结构体
// 前后端交互的 JSON 必须符合此结构
type MsgPacket struct {
	Type    int    `json:"type"`    // 消息类型
	FromID  string `json:"from_id"` // 发送者 UserID
	ToID    string `json:"to_id"`   // 接收者 UserID
	Content string `json:"content"` // 消息内容
	Seq     int64  `json:"seq"`     // 时间戳 (毫秒)，用于排序和去重
}
