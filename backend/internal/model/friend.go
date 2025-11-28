package model

import (
	"time"
)

type FriendRequest struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:id"`
	FromUserID    string    `gorm:"column:from_user_id;type:varchar(64);not null;index"`
	ToUserID      string    `gorm:"column:to_user_id;type:varchar(64);not null;index"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;type:varchar(255)"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	HandlerUserID string    `gorm:"column:handler_user_id;type:varchar(64)"`
	HandleMsg     string    `gorm:"column:handle_msg;type:varchar(255)"`
	HandledAt     time.Time `gorm:"column:handled_at"`
	Ex            string    `gorm:"column:ex;type:text"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

type Friend struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id"`
	OwnerUserID    string    `gorm:"column:owner_user_id;type:varchar(64);not null;index"`
	FriendUserID   string    `gorm:"column:friend_user_id;type:varchar(64);not null;index"`
	Remark         string    `gorm:"column:remark;type:varchar(255)"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(64)"`
	Ex             string    `gorm:"column:ex;type:text"`
	IsPinned       bool      `gorm:"column:is_pinned"`
}

// TableName 指定表名
func (Friend) TableName() string {
	return "friends"
}

type Black struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id"`
	OwnerUserID    string    `gorm:"column:owner_user_id;type:varchar(64);not null;index"`
	BlockUserID    string    `gorm:"column:block_user_id;type:varchar(64);not null;index"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(64)"`
	Ex             string    `gorm:"column:ex;type:text"`
}

func (Black) TableName() string {
	return "blacks"
}
