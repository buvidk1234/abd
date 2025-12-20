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

type User struct {
	// 基本信息
	UserID       int64  `gorm:"column:user_id;primaryKey;autoIncrement" json:"userID,string"`
	Username     string `gorm:"column:username;type:varchar(32);uniqueIndex;not null;comment:用户名" json:"username"`
	PasswordHash string `gorm:"column:password;type:varchar(255);not null;comment:加密后的密码" json:"-"`
	Phone        string `gorm:"column:phone;type:varchar(20);uniqueIndex;default:null;comment:手机号" json:"phone,omitempty"`
	Email        string `gorm:"column:email;type:varchar(64);index;default:null;comment:邮箱" json:"email,omitempty"`
	// 个人资料
	Nickname  string    `gorm:"column:nickname;type:varchar(64);default:'';comment:昵称" json:"nickname"`
	AvatarURL string    `gorm:"column:face_url;type:varchar(255);default:'';comment:头像链接" json:"faceURL"`
	Gender    int32     `gorm:"column:gender;type:tinyint;default:0;comment:性别，0：不愿透露，1：男，2：女" json:"gender"`
	Signature string    `gorm:"column:signature;type:varchar(255);default:'';comment:个性签名" json:"signature"`
	Birth     time.Time `gorm:"column:birth;type:date;default:null" json:"birth"`
	// 账号设置
	AppManagerLevel  int32 `gorm:"column:app_manager_level;type:tinyint;default:1;comment:管理级别" json:"appManagerLevel"`
	GlobalRecvMsgOpt int32 `gorm:"column:global_recv_msg_opt;type:tinyint;default:0;comment:全局免打扰0正常1免打扰" json:"globalRecvMsgOpt"`
	Status           int32 `gorm:"column:status;type:tinyint;default:1;comment:状态1正常2封禁" json:"status"`
	// 扩展字段
	Ex        string         `gorm:"column:ex;type:json;comment:扩展字段" json:"ex"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:注册时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string { return "users" }

/*func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == 0 {
		u.UserID = node.Generate().Int64()
	}
	return
}*/

func SelectUserInfo(tx *gorm.DB) *gorm.DB {
	return tx.Select("user_id", "username", "nickname", "face_url", "gender", "signature")
}

type UserEx struct {
	City     string `json:"city"`
	JobTitle string `json:"job_title"`
}

func (u *User) SetEx(exObj UserEx) error {
	bytes, err := json.Marshal(exObj)
	if err != nil {
		return err
	}
	u.Ex = string(bytes)
	return nil
}
func (u *User) GetEx() (*UserEx, error) {
	if u.Ex == "" {
		return &UserEx{}, nil
	}
	var exObj UserEx
	err := json.Unmarshal([]byte(u.Ex), &exObj)
	return &exObj, err
}
