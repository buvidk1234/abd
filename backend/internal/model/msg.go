package model

type SeqConversation struct {
	// ID ConversationID就是对象ID (GroupID 或 (UserAID,UserBID))
	ID string `gorm:"column:id;type:varchar(64);primaryKey"`

	// 序列号类型：1=UserSeq, 2=GroupSeq
	SeqType int32 `gorm:"column:seq_type;primaryKey"`

	// 当前最大值
	MaxSeq int64 `gorm:"column:max_seq;not null"`
}

func (SeqConversation) TableName() string {
	return "seq_conversations"
}

type Conversation struct {
	// 1. 联合主键
	OwnerID        string `gorm:"column:owner_id;type:varchar(64);primaryKey"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);primaryKey"`

	// 2. 会话类型
	ConvType int32 `gorm:"column:conv_type;not null"` // 1=单聊, 2=群聊

	// 3. 核心状态控制 (这就是你提到的 Status)
	// 1=正常, 2=已删除(隐藏), 3=被封禁/异常
	Status int32 `gorm:"column:status;default:1;index"`

	// 4. 个性化设置
	UnreadCount int32  `gorm:"column:unread_count;default:0"`
	IsPinned    bool   `gorm:"column:is_pinned;default:false;index"` // 置顶需索引，方便排序
	IsMuted     bool   `gorm:"column:is_muted;default:false"`
	ShowName    string `gorm:"column:show_name;type:varchar(128)"` // 备注名

	// 5. 同步位点 (Checkpoint)
	MinSeq  int64 `gorm:"column:min_seq"`  // 会话内最小 GroupSeq
	ReadSeq int64 `gorm:"column:read_seq"` // 已读到的 GroupSeq
	SyncSeq int64 `gorm:"column:sync_seq"` // 仅用于 Timeline 模式：最后一条同步到的 Timeline Seq,如果用户传的是0，同步最近min(100, SyncSeq-0)条消息。 单设备

	// 6. 冗余显示 (兜底用)
	LastMsgID       int64  `gorm:"column:last_msg_id"`
	LastMsgTime     int64  `gorm:"column:last_msg_time;index"` // 用于列表排序
	LastMsgSnapshot string `gorm:"column:last_msg_snapshot;type:varchar(500)"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// 多设备支持
type DeviceCheckpoint struct {
	UserID         string `gorm:"primaryKey"`
	DeviceID       string `gorm:"primaryKey"` // 重点：区分 iPhone, Windows, Android
	ConversationID string `gorm:"primaryKey"` // 如果是Timeline模式，这里不需要ConvID，只需要Seq

	SyncSeq int64 // 这个设备的同步进度
}

func (DeviceCheckpoint) TableName() string {
	return "device_checkpoints"
}

type Message struct {
	// 1. 核心ID：使用雪花算法(Snowflake)，不要用自增ID
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:false"`

	// 2. 归属与排序
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);index:idx_conv_seq,priority:1;not null"`
	Seq            int64  `gorm:"column:seq;index:idx_conv_seq,priority:2;not null"` // 联合唯一索引的核心，群内递增

	// 3. 消息本体
	SenderID    string `gorm:"column:sender_id;type:varchar(64);not null"`
	ClientMsgID string `gorm:"column:client_msg_id;type:varchar(64);index"` // 用于去重
	MsgType     int32  `gorm:"column:msg_type;not null"`                    // 1=文本, 2=图片...
	Content     string `gorm:"column:content;type:longtext"`                // 完整内容(JSON)

	// 4. 引用与状态
	RefMsgID int64 `gorm:"column:ref_msg_id"`       // 引用/回复的消息ID
	Status   int32 `gorm:"column:status;default:0"` // 0=正常, 1=撤回(撤回详情表),。删除(部分用户不可见有相关表)前端自行维护

	// 5. 时间
	SendTime   int64 `gorm:"column:send_time;index;not null"`            // 发送时间戳(ms)
	CreateTime int64 `gorm:"column:create_time;autoCreateTime;not null"` // 入库时间戳(ms)

	// 6. 冗余字段
	// 	SenderNickname   string `gorm:"column:sender_nickname;type:varchar(128)"`
	// 	SenderAvatarURL  string `gorm:"column:sender_avatar_url;type:varchar(255)"`

	ConvType int32  `gorm:"column:conv_type;not null"`
	TargetID string `gorm:"column:target_id;type:varchar(64);not null"`

	// 分库分表策略：通常按 ConversationID 取模分表
}

func (Message) TableName() string {
	return "messages"
}

type SeqUser struct {
	ID     string `gorm:"column:id;type:varchar(64);primaryKey"`
	MaxSeq int64  `gorm:"column:max_seq;not null"`
}

func (SeqUser) TableName() string {
	return "seq_users"
}

type UserTimeline struct {
	// 1. 联合主键 (聚簇索引)：查询 WHERE owner_id=? AND seq>?
	OwnerID string `gorm:"column:owner_id;type:varchar(64);primaryKey;index:idx_owner_seq,priority:1"`
	Seq     int64  `gorm:"column:seq;primaryKey;index:idx_owner_seq,priority:2"` // 用户维度的全局Seq

	// 2. 来源定位
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);not null"`
	MsgID          int64  `gorm:"column:msg_id;not null"` // 关联 messages.id
	RefMsgSeq      int64  `gorm:"column:ref_msg_seq"`     // 冗余一份 GroupSeq，方便客户端校验补洞

	// 3. 列表页快照 (List View)
	MsgType  int32  `gorm:"column:msg_type;not null"`
	SenderID string `gorm:"column:sender_id;type:varchar(64)"`
	Snapshot string `gorm:"column:snapshot;type:varchar(1000)"` // 摘要："[图片]", "你好..."

	// 4. 时间与状态
	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`

	// 分库分表策略：按 OwnerID 取模分表
}

func (UserTimeline) TableName() string {
	return "user_timelines"
}

// 1对1
type MsgRevoke struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;column:id"`
	MsgID    string `gorm:"column:msg_id;type:varchar(128);index"`
	Role     int32  `gorm:"column:role"`
	UserID   string `gorm:"column:user_id;type:varchar(64);index"`
	Nickname string `gorm:"column:nickname;type:varchar(128)"`
	Time     int64  `gorm:"column:time;index"`
}

func (MsgRevoke) TableName() string {
	return "msg_revokes"
}

// 1对多
// type MsgAt struct {
// 	ID       uint   `gorm:"primaryKey;autoIncrement;column:id"`
// 	MsgID    string `gorm:"column:msg_id;type:varchar(128);index"`
// 	AtUserID string `gorm:"column:at_user_id;type:varchar(64);index"`
// }

// func (MsgAt) TableName() string {
// 	return "msg_ats"
// }

// 1对多
// type MsgDelete struct {
// 	ID     uint   `gorm:"primaryKey;autoIncrement;column:id"`
// 	MsgID  string `gorm:"column:msg_id;type:varchar(128);index"`
// 	UserID string `gorm:"column:user_id;type:varchar(64);index"`
// 	Time   int64  `gorm:"column:time;index"`
// }

// func (MsgDelete) TableName() string {
// 	return "msg_deletes"
// }
