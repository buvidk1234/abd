package model

type SeqConversation struct {
	// ID ConversationID就是对象ID (GroupID 或 (UserAID,UserBID))
	ID string `gorm:"column:id;type:varchar(64);primaryKey" json:"id"`

	// 序列号类型：1=UserSeq, 2=GroupSeq
	SeqType int32 `gorm:"column:seq_type;primaryKey" json:"seq_type"`

	// 当前最大值
	MaxSeq int64 `gorm:"column:max_seq;not null" json:"max_seq"`
	// 当前最小值
	MinSeq int64 `gorm:"column:min_seq;not null" json:"min_seq"`
}

func (SeqConversation) TableName() string {
	return "seq_conversations"
}

type Conversation struct {
	// 1. 联合主键
	OwnerID        int64  `gorm:"column:owner_id;primaryKey" json:"owner_id,string"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);primaryKey" json:"conversation_id"`

	// 2. 会话类型
	ConvType int32 `gorm:"column:conv_type;not null" json:"conv_type"` // 1=单聊, 2=群聊

	// 3. 核心状态控制 (这就是你提到的 Status)
	// 1=正常, 2=已删除(隐藏), 3=被封禁/异常
	Status int32 `gorm:"column:status;default:1;index" json:"status"`

	// 4. 个性化设置
	UnreadCount int32  `gorm:"column:unread_count;default:0" json:"unread_count"`
	IsPinned    bool   `gorm:"column:is_pinned;default:false;index" json:"is_pinned"` // 置顶需索引，方便排序
	IsMuted     bool   `gorm:"column:is_muted;default:false" json:"is_muted"`
	ShowName    string `gorm:"column:show_name;type:varchar(128)" json:"show_name"` // 备注名

	// 5. 同步位点 (Checkpoint)
	MinSeq  int64 `gorm:"column:min_seq" json:"min_seq"`   // 会话内最小 GroupSeq
	MaxSeq  int64 `gorm:"column:max_seq" json:"max_seq"`   // 能检索的最大值
	ReadSeq int64 `gorm:"column:read_seq" json:"read_seq"` // 已读到的 GroupSeq
	SyncSeq int64 `gorm:"column:sync_seq" json:"sync_seq"` // 仅用于 Timeline 模式：最后一条同步到的 Timeline Seq,如果用户传的是0，同步最近min(100, SyncSeq-0)条消息。 单设备

	// 6. 冗余显示 (兜底用)
	LastMsgID       int64  `gorm:"column:last_msg_id" json:"last_msg_id,string"`
	LastMsgTime     int64  `gorm:"column:last_msg_time;index" json:"last_msg_time"` // 用于列表排序
	LastMsgSnapshot string `gorm:"column:last_msg_snapshot;type:varchar(500)" json:"last_msg_snapshot"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// 多设备支持
type DeviceCheckpoint struct {
	UserID         int64  `gorm:"primaryKey" json:"user_id,string"`
	DeviceID       string `gorm:"primaryKey" json:"device_id"`       // 重点：区分 iPhone, Windows, Android
	ConversationID string `gorm:"primaryKey" json:"conversation_id"` // 如果是Timeline模式，这里不需要ConvID，只需要Seq

	SyncSeq int64 `json:"sync_seq"` // 这个设备的同步进度
}

func (DeviceCheckpoint) TableName() string {
	return "device_checkpoints"
}

type Message struct {
	// 1. 核心ID：使用雪花算法(Snowflake)，不要用自增ID
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:false" json:"id,string"`

	// 2. 归属与排序
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);index:idx_conv_seq,priority:1;not null" json:"conversation_id"`
	Seq            int64  `gorm:"column:seq;index:idx_conv_seq,priority:2;not null" json:"seq"` // 联合唯一索引的核心，群内递增

	// 3. 消息本体
	SenderID    int64  `gorm:"column:sender_id;not null" json:"sender_id,string"`
	ClientMsgID string `gorm:"column:client_msg_id;type:varchar(64);index" json:"client_msg_id"` // 用于去重
	MsgType     int32  `gorm:"column:msg_type;not null" json:"msg_type"`                         // 1=文本, 2=图片...
	Content     string `gorm:"column:content;type:longtext" json:"content"`                      // 完整内容(JSON)

	// 4. 引用与状态
	RefMsgID int64 `gorm:"column:ref_msg_id" json:"ref_msg_id,string"` // 引用/回复的消息ID
	Status   int32 `gorm:"column:status;default:0" json:"status"`      // 0=正常, 1=撤回(撤回详情表),。删除(部分用户不可见有相关表)前端自行维护

	// 5. 时间
	SendTime   int64 `gorm:"column:send_time;index;not null" json:"send_time"`              // 发送时间戳(ms)
	CreateTime int64 `gorm:"column:create_time;autoCreateTime;not null" json:"create_time"` // 入库时间戳(ms)

	// 6. 冗余字段
	// 	SenderNickname   string `gorm:"column:sender_nickname;type:varchar(128)"`
	// 	SenderAvatarURL  string `gorm:"column:sender_avatar_url;type:varchar(255)"`

	ConvType int32 `gorm:"column:conv_type;not null" json:"conv_type"`
	TargetID int64 `gorm:"column:target_id;not null" json:"target_id,string"`

	// 分库分表策略：通常按 ConversationID 取模分表
}

func (Message) TableName() string {
	return "messages"
}

type SeqUser struct {
	UserID         int64  `gorm:"column:user_id;primaryKey" json:"user_id,string"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);primaryKey" json:"conversation_id"`

	MinSeq  int64 `gorm:"column:min_seq;not null" json:"min_seq"`   // 会话内最小 UserSeq
	ReadSeq int64 `gorm:"column:read_seq;not null" json:"read_seq"` // 已读到的 UserSeq
	MaxSeq  int64 `gorm:"column:max_seq;not null" json:"max_seq"`   // 当前最大值
}

func (SeqUser) TableName() string {
	return "seq_users"
}

// 用来实现feed流的用户时间线，对于单聊和小群聊不需要,用来实现发帖，朋友圈等功能
type UserTimeline struct {
	// 1. 联合主键 (聚簇索引)：查询 WHERE owner_id=? AND seq>?
	OwnerID int64 `gorm:"column:owner_id;primaryKey;index:idx_owner_seq,priority:1"`
	Seq     int64 `gorm:"column:seq;primaryKey;index:idx_owner_seq,priority:2"` // 用户维度的全局Seq

	// 2. 来源定位
	ConversationID string `gorm:"column:conversation_id;type:varchar(64);not null"`
	MsgID          int64  `gorm:"column:msg_id;not null"` // 关联 messages.id
	RefMsgSeq      int64  `gorm:"column:ref_msg_seq"`     // 冗余一份 GroupSeq，方便客户端校验补洞

	// 3. 列表页快照 (List View)
	MsgType  int32  `gorm:"column:msg_type;not null"`
	SenderID int64  `gorm:"column:sender_id"`
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
	UserID   int64  `gorm:"column:user_id;index"`
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
