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
	ID uint `gorm:"primaryKey;autoIncrement;column:id"`
	//GroupID               string    `gorm:"column:group_id;type:varchar(64);uniqueIndex;not null"`
	GroupName             string    `gorm:"column:group_name;type:varchar(128);not null"`
	Notification          string    `gorm:"column:notification;type:text;comment:群公告"`
	Introduction          string    `gorm:"column:introduction;type:text;comment:群简介"`
	AvatarURL             string    `gorm:"column:avatar_url;type:varchar(255);comment:群头像"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
	Ex                    string    `gorm:"column:ex;type:text;comment:扩展字段"`
	Status                int32     `gorm:"column:status;comment:群状态"`
	CreatorUserID         string    `gorm:"column:creator_user_id;type:varchar(64);comment:群创建者ID"`
	GroupType             int32     `gorm:"column:group_type;comment:群类型"`
	NeedVerification      int32     `gorm:"column:need_verification;comment:是否需要验证"`
	LookMemberInfo        int32     `gorm:"column:look_member_info;comment:是否允许查看成员信息"`
	ApplyMemberFriend     int32     `gorm:"column:apply_member_friend;comment:是否允许申请加群成员为好友"`
	NotificationUpdatedAt time.Time `gorm:"column:notification_updated_at;comment:群公告最后更新时间"`
	NotificationUserID    string    `gorm:"column:notification_user_id;type:varchar(64);comment:最后更新群公告的用户ID"`
}

func (Group) TableName() string {
	return "groups"
}

// 群成员表
type GroupMember struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id;comment:自增ID"`
	GroupID        string    `gorm:"column:group_id;type:varchar(64);index;not null;comment:群ID"`
	UserID         string    `gorm:"column:user_id;type:varchar(64);index;not null;comment:用户ID"`
	Nickname       string    `gorm:"column:nickname;type:varchar(128);comment:昵称"`
	AvatarURL      string    `gorm:"column:avatar_url;type:varchar(255);comment:头像"`
	RoleLevel      int32     `gorm:"column:role_level;comment:角色等级"`
	JoinedAt       time.Time `gorm:"column:joined_at;autoCreateTime"`
	JoinSource     int32     `gorm:"column:join_source;comment:加入来源"`
	InviterUserID  string    `gorm:"column:inviter_user_id;type:varchar(64);comment:邀请人ID"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(64);comment:操作人ID"`
	MuteEndAt      time.Time `gorm:"column:mute_end_at;comment:禁言结束时间"`
	Ex             string    `gorm:"column:ex;type:text;comment:扩展字段"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
