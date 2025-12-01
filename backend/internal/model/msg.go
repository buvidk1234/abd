package model

type SeqConversation struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;column:id"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(128);index"`
	MinSeq         int64  `gorm:"column:min_seq;index"`
	MaxSeq         int64  `gorm:"column:max_seq;index"`
}

func (SeqConversation) TableName() string {
	return "seq_conversations"
}

type Conversation struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;column:id"`
	UserID         string `gorm:"column:user_id;type:varchar(64);index"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(128);index"`
	MinSeq         int64  `gorm:"column:min_seq;index"`
	MaxSeq         int64  `gorm:"column:max_seq;index"`
	ReadSeq        int64  `gorm:"column:read_seq;index"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type Message struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;column:id"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(128);index"`
	SenderID       string `gorm:"column:sender_id;type:varchar(64);index"`
	ContentType    int32  `gorm:"column:content_type"`
	Content        string `gorm:"column:content;type:text"`
	Seq            int64  `gorm:"column:seq;index"`
	CreatedAt      int64  `gorm:"column:created_at;autoCreateTime;index"`
}

// type Message struct {
// 	ID          uint   `gorm:"primaryKey;autoIncrement;column:id"`
// 	SendID      string `gorm:"column:send_id;type:varchar(64);index"`
// 	SessionType int32  `gorm:"column:session_type;index"`
// 	RecvID      string `gorm:"column:recv_id;type:varchar(64);index"`
// 	GroupID     string `gorm:"column:group_id;type:varchar(64);index"`

// 	SenderPlatformID int32  `gorm:"column:sender_platform_id"`
// 	SenderNickname   string `gorm:"column:sender_nickname;type:varchar(128)"`
// 	SenderAvatarURL  string `gorm:"column:sender_avatar_url;type:varchar(255)"`

// 	ClientMsgID string `gorm:"column:client_msg_id;type:varchar(128);index"`
// 	ServerMsgID string `gorm:"column:server_msg_id;type:varchar(128);index"`
// 	Seq         int64  `gorm:"column:seq;index"`
// 	SentAt      int64  `gorm:"column:sent_at;index"`
// 	CreatedAt   int64  `gorm:"column:created_at;index"`

// 	MsgFrom     int32  `gorm:"column:msg_from"`
// 	ContentType int32  `gorm:"column:content_type"`
// 	Content     string `gorm:"column:content;type:text"`

// 	Status int32 `gorm:"column:status;index"`

// 	ConversationID string `gorm:"column:conversation_id;type:varchar(128);index"` //
// }

// 1对1
// type MsgRevoke struct {
// 	ID       uint   `gorm:"primaryKey;autoIncrement;column:id"`
// 	MsgID    string `gorm:"column:msg_id;type:varchar(128);index"`
// 	Role     int32  `gorm:"column:role"`
// 	UserID   string `gorm:"column:user_id;type:varchar(64);index"`
// 	Nickname string `gorm:"column:nickname;type:varchar(128)"`
// 	Time     int64  `gorm:"column:time;index"`
// }

// func (MsgRevoke) TableName() string {
// 	return "msg_revokes"
// }

// // 1对多
// type MsgAt struct {
// 	ID       uint   `gorm:"primaryKey;autoIncrement;column:id"`
// 	MsgID    string `gorm:"column:msg_id;type:varchar(128);index"`
// 	AtUserID string `gorm:"column:at_user_id;type:varchar(64);index"`
// }

// func (MsgAt) TableName() string {
// 	return "msg_ats"
// }

// // 1对多
// type MsgDelete struct {
// 	ID     uint   `gorm:"primaryKey;autoIncrement;column:id"`
// 	MsgID  string `gorm:"column:msg_id;type:varchar(128);index"`
// 	UserID string `gorm:"column:user_id;type:varchar(64);index"`
// 	Time   int64  `gorm:"column:time;index"`
// }

// func (MsgDelete) TableName() string {
// 	return "msg_deletes"
// }

// type MsgOfflinePush struct {
// 	Title         string
// 	Desc          string
// 	IOSPushSound  string
// 	IOSBadgeCount bool
// }

// func (Message) TableName() string {
// 	return "messages"
// }
