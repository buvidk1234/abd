package model

import (
	"time"
)

// 群申请表
type GroupRequest struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:id"`
	UserID        string    `gorm:"column:user_id;type:varchar(64);index;not null"`
	GroupID       string    `gorm:"column:group_id;type:varchar(64);index;not null"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;type:text"`
	HandledMsg    string    `gorm:"column:handled_msg;type:text"`
	RequestedAt   time.Time `gorm:"column:requested_at;autoCreateTime"`
	HandleUserID  string    `gorm:"column:handle_user_id;type:varchar(64)"`
	HandledAt     time.Time `gorm:"column:handled_at"`
	JoinSource    int32     `gorm:"column:join_source"`
	InviterUserID string    `gorm:"column:inviter_user_id;type:varchar(64)"`
	Ex            string    `gorm:"column:ex;type:text"`
}

func (GroupRequest) TableName() string {
	return "group_requests"
}

// 群组表
type Group struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement;column:id"`
	GroupID               string    `gorm:"column:group_id;type:varchar(64);uniqueIndex;not null"`
	GroupName             string    `gorm:"column:group_name;type:varchar(128);not null"`
	Notification          string    `gorm:"column:notification;type:text"`
	Introduction          string    `gorm:"column:introduction;type:text"`
	AvatarURL             string    `gorm:"column:avatar_url;type:varchar(255)"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
	Ex                    string    `gorm:"column:ex;type:text"`
	Status                int32     `gorm:"column:status"`
	CreatorUserID         string    `gorm:"column:creator_user_id;type:varchar(64)"`
	GroupType             int32     `gorm:"column:group_type"`
	NeedVerification      int32     `gorm:"column:need_verification"`
	LookMemberInfo        int32     `gorm:"column:look_member_info"`
	ApplyMemberFriend     int32     `gorm:"column:apply_member_friend"`
	NotificationUpdatedAt time.Time `gorm:"column:notification_updated_at"`
	NotificationUserID    string    `gorm:"column:notification_user_id;type:varchar(64)"`
}

func (Group) TableName() string {
	return "groups"
}

// 群成员表
type GroupMember struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id"`
	GroupID        string    `gorm:"column:group_id;type:varchar(64);index;not null"`
	UserID         string    `gorm:"column:user_id;type:varchar(64);index;not null"`
	Nickname       string    `gorm:"column:nickname;type:varchar(128)"`
	AvatarURL      string    `gorm:"column:avatar_url;type:varchar(255)"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinedAt       time.Time `gorm:"column:joined_at;autoCreateTime"`
	JoinSource     int32     `gorm:"column:join_source"`
	InviterUserID  string    `gorm:"column:inviter_user_id;type:varchar(64)"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(64)"`
	MuteEndAt      time.Time `gorm:"column:mute_end_at"`
	Ex             string    `gorm:"column:ex;type:text"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
