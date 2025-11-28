package model

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

// 定义一个全局的 Node，用于生成 ID
var node *snowflake.Node

func init() {
	var err error
	// 1 表示节点 ID，分布式部署时不同机器要不一样
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

// User 用户表
type User struct {
	// ================= 核心身份信息 =================

	// UserID 使用 String 类型 (推荐 UUID 或 雪花算法)，长度设为 64 足够容纳大多数 ID
	UserID string `gorm:"column:user_id;primaryKey;type:varchar(64)" json:"userID"`

	// 用户名，唯一索引，用于登录
	Username string `gorm:"column:username;type:varchar(32);uniqueIndex;not null;comment:用户名" json:"username"`

	// 密码哈希，注意 json:"-" 确保 API 不会返回密码字段
	PasswordHash string `gorm:"column:password;type:varchar(255);not null;comment:加密后的密码" json:"-"`

	// 手机号，唯一索引 (IM 系统常用手机号登录)
	Phone string `gorm:"column:phone;type:varchar(20);uniqueIndex;default:null;comment:手机号" json:"phone,omitempty"`

	// 邮箱 (可选)
	Email string `gorm:"column:email;type:varchar(64);index;default:null;comment:邮箱" json:"email,omitempty"`

	// ================= 个人资料 (Profile) =================

	// 昵称
	Nickname string `gorm:"column:nickname;type:varchar(64);default:'';comment:昵称" json:"nickname"`

	// 头像 URL (原 face_url)
	AvatarURL string `gorm:"column:face_url;type:varchar(255);default:'';comment:头像链接" json:"faceURL"`

	// 性别 0:保密 1:男 2:女
	Gender int32 `gorm:"column:gender;type:tinyint;default:0;comment:性别" json:"gender"`

	// 个性签名
	Signature string `gorm:"column:signature;type:varchar(255);default:'';comment:个性签名" json:"signature"`

	// 生日
	Birth time.Time `gorm:"column:birth;type:date;default:null" json:"birth"`

	// ================= 状态与权限 =================

	// 管理员级别 1:普通用户 2:App管理员 (修正了拼写 AppMangerLevel -> AppManagerLevel)
	AppManagerLevel int32 `gorm:"column:app_manager_level;type:tinyint;default:1;comment:管理级别" json:"appManagerLevel"`

	// 全局免打扰设置 0:正常接收 1:全局免打扰 (GlobalReceiveMessageOption)
	GlobalRecvMsgOpt int32 `gorm:"column:global_recv_msg_opt;type:tinyint;default:0;comment:全局免打扰0正常1免打扰" json:"globalRecvMsgOpt"`

	// 账号状态 1:正常 2:封禁
	Status int32 `gorm:"column:status;type:tinyint;default:1;comment:状态1正常2封禁" json:"status"`

	// ================= 扩展与系统字段 =================

	// 扩展字段 (建议存储 JSON 字符串)，用于存放不确定的配置信息
	Ex string `gorm:"column:ex;type:json;comment:扩展字段" json:"ex"`

	// 创建时间 (GORM 自动维护)
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime;comment:注册时间" json:"createTime"`

	// 更新时间 (GORM 自动维护)
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`

	// 软删除 (删除用户时，不会物理删除，而是更新该字段)
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// ================= 辅助方法 =================

// BeforeCreate GORM 钩子
// 在插入数据库之前，自动生成有序的 UserID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == "" {
		// 生成 int64 的 ID 并转为 string
		u.UserID = node.Generate().String()
	}
	return
}

// ExObject 用于解析 Ex 字段的辅助结构 (示例)
type UserEx struct {
	City     string `json:"city"`
	JobTitle string `json:"job_title"`
}

// SetEx 设置扩展字段
func (u *User) SetEx(exObj UserEx) error {
	bytes, err := json.Marshal(exObj)
	if err != nil {
		return err
	}
	u.Ex = string(bytes)
	return nil
}

// GetEx 获取扩展字段
func (u *User) GetEx() (*UserEx, error) {
	if u.Ex == "" {
		return &UserEx{}, nil
	}
	var exObj UserEx
	err := json.Unmarshal([]byte(u.Ex), &exObj)
	return &exObj, err
}
