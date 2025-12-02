package constant

const (
	// ConvType.
	SingleChatType       = 1
	GroupChatType        = 2
	NotificationChatType = 3
)

const (
	// --- 基础内容 (确实是 ContentType) ---
	MsgTypeText  = 101
	MsgTypeImage = 102
	MsgTypeVideo = 103

	// --- 业务信令 (这叫 ContentType 就不合适了) ---
	MsgTypeRevoke     = 201 // 撤回
	MsgTypeReadReport = 202 // 已读回执
	MsgTypeTyping     = 203 // "正在输入中..."

	// --- 群组事件 (这也是业务逻辑) ---
	MsgTypeMemberJoin = 301 // "张三加入群聊"
	MsgTypeGroupMute  = 302 // "群主开启了全员禁言"
)
