package model

import (
	"time"

	"gorm.io/gorm"
)

type FriendRequest struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;column:id;comment:主键ID" json:"id,string"`
	FromUserID    int64     `gorm:"column:from_user_id;not null;index;comment:发送方" json:"fromUserID,string"`
	ToUserID      int64     `gorm:"column:to_user_id;not null;index;comment:接收方" json:"toUserID,string"`
	HandleResult  int32     `gorm:"column:handle_result;type:tinyint;default:0;comment:处理结果，1：同意，2：拒绝，0：未处理" json:"handleResult"`
	ReqMsg        string    `gorm:"column:req_msg;type:varchar(255);comment:请求消息" json:"reqMsg"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	HandlerUserID int64     `gorm:"column:handler_user_id;comment:处理方" json:"handlerUserID,string"`
	HandleMsg     string    `gorm:"column:handle_msg;type:varchar(255);comment:处理消息" json:"handleMsg"`
	HandledAt     time.Time `gorm:"column:handled_at" json:"handledAt"`
	Ex            string    `gorm:"column:ex;type:text;comment:扩展字段" json:"ex"`
	FromUser      User      `gorm:"foreignKey:FromUserID;references:UserID"`
	ToUser        User      `gorm:"foreignKey:ToUserID;references:UserID"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

func SelectFriendRequestInfo(tx *gorm.DB) *gorm.DB {
	return tx.Select(
		"id",
		"from_user_id",
		"to_user_id",
		"handle_result",
		"req_msg",
		"created_at",
		"updated_at",
		"handle_msg",
		"handled_at",
	)
}

type Friend struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id,string"`
	OwnerUserID    int64     `gorm:"column:owner_user_id;not null;index;comment:所有者" json:"ownerUserID,string"`
	FriendUserID   int64     `gorm:"column:friend_user_id;not null;index;comment:好友" json:"friendUserID,string"`
	Remark         string    `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	AddSource      int32     `gorm:"column:add_source;type:tinyint;default:0;comment:添加来源，1：搜索添加，2：扫码添加，3：邀请添加，4：推荐添加，5：其他" json:"addSource"`
	OperatorUserID int64     `gorm:"column:operator_user_id;comment:操作人" json:"operatorUserID,string"`
	Ex             string    `gorm:"column:ex;type:text" json:"ex"`
	IsPinned       bool      `gorm:"column:is_pinned;default:false;comment:是否置顶" json:"isPinned"`
	OwnerUser      User      `gorm:"foreignKey:OwnerUserID;references:UserID"`
	FriendUser     User      `gorm:"foreignKey:FriendUserID;references:UserID"`
}

func SelectFriendInfo(tx *gorm.DB) *gorm.DB {
	return tx.Select("id", "owner_user_id", "friend_user_id", "remark", "created_at", "add_source", "is_pinned")
}

// TableName 指定表名
func (Friend) TableName() string {
	return "friends"
}

type Black struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id,string"`
	OwnerUserID    int64     `gorm:"column:owner_user_id;not null;index;comment:所有者" json:"ownerUserID,string"`
	BlockUserID    int64     `gorm:"column:block_user_id;not null;index;comment:黑名单用户" json:"blockUserID,string"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	AddSource      int32     `gorm:"column:add_source;type:tinyint;default:0;comment:添加来源，1：搜索添加，2：扫码添加，3：邀请添加，4：推荐添加，5：其他" json:"addSource"`
	OperatorUserID int64     `gorm:"column:operator_user_id;comment:操作人" json:"operatorUserID,string"`
	Ex             string    `gorm:"column:ex;type:text" json:"ex"`
}

func (Black) TableName() string {
	return "blacks"
}
