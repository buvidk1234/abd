package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarURL"`
}

func (u *UserService) UserRegister(ctx context.Context, req UserRegisterReq) (string, error) {
	// 检查用户名是否已存在
	var count int64
	if err := u.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return "", err
	}
	if count > 0 {
		return "", gorm.ErrDuplicatedKey
	}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Phone:        req.Phone,
		Email:        req.Email,
		Nickname:     req.Nickname,
		AvatarURL:    req.AvatarURL,
	}
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return "", err
	}
	return user.UserID, nil
}

// UpdateUserInfoReq 用户信息更新请求
type UpdateUserInfoReq struct {
	UserID    string  `json:"userID" binding:"required"`
	Nickname  string  `json:"nickname"`
	AvatarURL string  `json:"avatarURL"`
	Gender    int32   `json:"gender"`
	Signature string  `json:"signature"`
	Birth     *string `json:"birth"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Ex        string  `json:"ex"`
}

// UpdateUserInfo
func (u *UserService) UpdateUserInfo(ctx context.Context, req UpdateUserInfoReq) error {
	var user model.User
	if err := u.db.WithContext(ctx).First(&user, "user_id = ?", req.UserID).Error; err != nil {
		return err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Gender != 0 {
		user.Gender = req.Gender
	}
	if req.Signature != "" {
		user.Signature = req.Signature
	}
	if req.Birth != nil && *req.Birth != "" {
		t, err := time.Parse("2006-01-02", *req.Birth)
		if err == nil {
			user.Birth = t
		}
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Ex != "" {
		user.Ex = req.Ex
	}
	return u.db.WithContext(ctx).Save(&user).Error
}

type GetUsersPublicInfoReq struct {
	UserIDs []string `json:"userIDs"`
}

func (u *UserService) GetUsersPublicInfo(ctx context.Context, req GetUsersPublicInfoReq) ([]dto.UserInfo, error) {
	var users []model.User
	if len(req.UserIDs) == 0 {
		u.db.WithContext(ctx).Find(&users)
	} else {
		if err := u.db.WithContext(ctx).Find(&users, "user_id IN ?", req.UserIDs).Error; err != nil {
			return nil, err
		}
	}
	var userInfos []dto.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, dto.ConvertToUserInfo(user))
	}
	return userInfos, nil
}

type UserLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *UserService) UserLogin(ctx context.Context, req UserLoginReq) error {
	var user model.User
	if err := u.db.WithContext(ctx).First(&user, "username = ?", req.Username).Error; err != nil {
		return errors.New("用户不存在")
	}
	//使用加密库比对密码
	bytePassword := []byte(req.Password)
	byteHashedPassword := []byte(user.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
