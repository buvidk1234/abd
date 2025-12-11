package model

import (
	"time"
)

// 群申请表
type GroupRequest struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID        int64     `gorm:"column:user_id;index;not null" json:"userID,string"`
	GroupID       string    `gorm:"column:group_id;type:varchar(64);index;not null" json:"groupID"`
	HandleResult  int32     `gorm:"column:handle_result" json:"handleResult"`
	ReqMsg        string    `gorm:"column:req_msg;type:text" json:"reqMsg"`
	HandledMsg    string    `gorm:"column:handled_msg;type:text" json:"handledMsg"`
	RequestedAt   time.Time `gorm:"column:requested_at;autoCreateTime" json:"requestedAt"`
	HandleUserID  int64     `gorm:"column:handle_user_id" json:"handleUserID,string"`
	HandledAt     time.Time `gorm:"column:handled_at" json:"handledAt"`
	JoinSource    int32     `gorm:"column:join_source" json:"joinSource"`
	InviterUserID int64     `gorm:"column:inviter_user_id" json:"inviterUserID,string"`
	Ex            string    `gorm:"column:ex;type:text" json:"ex"`
}

func (GroupRequest) TableName() string {
	return "group_requests"
}

// 群组表
type Group struct {
	ID uint `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	//GroupID               string    `gorm:"column:group_id;type:varchar(64);uniqueIndex;not null"`
	GroupName             string    `gorm:"column:group_name;type:varchar(128);not null" json:"groupName"`
	Notification          string    `gorm:"column:notification;type:text;comment:群公告" json:"notification"`
	Introduction          string    `gorm:"column:introduction;type:text;comment:群简介" json:"introduction"`
	AvatarURL             string    `gorm:"column:avatar_url;type:varchar(255);comment:群头像" json:"avatarURL"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Ex                    string    `gorm:"column:ex;type:text;comment:扩展字段" json:"ex"`
	Status                int32     `gorm:"column:status;comment:群状态" json:"status"`
	CreatorUserID         int64     `gorm:"column:creator_user_id;comment:群创建者ID" json:"creatorUserID,string"`
	GroupType             int32     `gorm:"column:group_type;comment:群类型" json:"groupType"`
	NeedVerification      int32     `gorm:"column:need_verification;comment:是否需要验证" json:"needVerification"`
	LookMemberInfo        int32     `gorm:"column:look_member_info;comment:是否允许查看成员信息" json:"lookMemberInfo"`
	ApplyMemberFriend     int32     `gorm:"column:apply_member_friend;comment:是否允许申请加群成员为好友" json:"applyMemberFriend"`
	NotificationUpdatedAt time.Time `gorm:"column:notification_updated_at;comment:群公告最后更新时间" json:"notificationUpdatedAt"`
	NotificationUserID    int64     `gorm:"column:notification_user_id;comment:最后更新群公告的用户ID" json:"notificationUserID,string"`
}

func (Group) TableName() string {
	return "groups"
}

// 群成员表
type GroupMember struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id;comment:自增ID" json:"id"`
	GroupID        string    `gorm:"column:group_id;type:varchar(64);index;not null;comment:群ID" json:"groupID"`
	UserID         int64     `gorm:"column:user_id;index;not null;comment:用户ID" json:"userID,string"`
	Nickname       string    `gorm:"column:nickname;type:varchar(128);comment:昵称" json:"nickname"`
	AvatarURL      string    `gorm:"column:avatar_url;type:varchar(255);comment:头像" json:"avatarURL"`
	RoleLevel      int32     `gorm:"column:role_level;comment:角色等级" json:"roleLevel"`
	JoinedAt       time.Time `gorm:"column:joined_at;autoCreateTime" json:"joinedAt"`
	JoinSource     int32     `gorm:"column:join_source;comment:加入来源" json:"joinSource"`
	InviterUserID  int64     `gorm:"column:inviter_user_id;comment:邀请人ID" json:"inviterUserID,string"`
	OperatorUserID int64     `gorm:"column:operator_user_id;comment:操作人ID" json:"operatorUserID,string"`
	MuteEndAt      time.Time `gorm:"column:mute_end_at;comment:禁言结束时间" json:"muteEndAt"`
	Ex             string    `gorm:"column:ex;type:text;comment:扩展字段" json:"ex"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
