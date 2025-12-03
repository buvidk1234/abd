package model

import (
	"time"
)

type FriendRequest struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;column:id;comment:主键ID"`
	FromUserID    int64     `gorm:"column:from_user_id;not null;index;comment:发送方"`
	ToUserID      int64     `gorm:"column:to_user_id;not null;index;comment:接收方"`
	HandleResult  int32     `gorm:"column:handle_result;type:tinyint;default:0;comment:处理结果，1：同意，2：拒绝，0：未处理"`
	ReqMsg        string    `gorm:"column:req_msg;type:varchar(255);comment:请求消息"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	HandlerUserID int64     `gorm:"column:handler_user_id;comment:处理方"`
	HandleMsg     string    `gorm:"column:handle_msg;type:varchar(255);comment:处理消息"`
	HandledAt     time.Time `gorm:"column:handled_at"`
	Ex            string    `gorm:"column:ex;type:text;comment:扩展字段"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

type Friend struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id"`
	OwnerUserID    int64     `gorm:"column:owner_user_id;not null;index;comment:所有者"`
	FriendUserID   int64     `gorm:"column:friend_user_id;not null;index;comment:好友"`
	Remark         string    `gorm:"column:remark;type:varchar(255);comment:备注"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	AddSource      int32     `gorm:"column:add_source;type:tinyint;default:0;comment:添加来源，1：搜索添加，2：扫码添加，3：邀请添加，4：推荐添加，5：其他"`
	OperatorUserID int64     `gorm:"column:operator_user_id;comment:操作人"`
	Ex             string    `gorm:"column:ex;type:text"`
	IsPinned       bool      `gorm:"column:is_pinned;default:false;comment:是否置顶"`
}

// TableName 指定表名
func (Friend) TableName() string {
	return "friends"
}

type Black struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id"`
	OwnerUserID    int64     `gorm:"column:owner_user_id;not null;index;comment:所有者"`
	BlockUserID    int64     `gorm:"column:block_user_id;not null;index;comment:黑名单用户"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	AddSource      int32     `gorm:"column:add_source;type:tinyint;default:0;comment:添加来源，1：搜索添加，2：扫码添加，3：邀请添加，4：推荐添加，5：其他"`
	OperatorUserID int64     `gorm:"column:operator_user_id;comment:操作人"`
	Ex             string    `gorm:"column:ex;type:text"`
}

func (Black) TableName() string {
	return "blacks"
}
