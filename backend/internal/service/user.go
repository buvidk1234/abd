package service

import (
	"backend/internal/api/apiresp/errs"
	"backend/internal/dto"
	"backend/internal/model"
	"backend/pkg/util"
	"context"
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

// SearchUserReq 搜索用户请求
type SearchUserReq struct {
	Keyword string `form:"keyword" binding:"required"`
}

func (u *UserService) SearchUser(ctx context.Context, req SearchUserReq) ([]dto.UserInfo, error) {
	var users []model.User
	// Search by username or nickname or phone
	// Note: checking username and nickname for now.
	if err := u.db.WithContext(ctx).Model(&model.User{}).Where("username LIKE ? OR nickname LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%").Find(&users).Error; err != nil {
		return nil, err
	}

	dtos := make([]dto.UserInfo, len(users))
	for i, user := range users {
		dtos[i] = dto.ConvertToUserInfo(user)
	}
	return dtos, nil
}

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

func (u *UserService) UserRegister(ctx context.Context, req UserRegisterReq) (int64, error) {
	// 检查用户名是否已存在
	var count int64
	if err := u.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errs.ErrUserExists
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
		return 0, err
	}
	return user.UserID, nil
}

// UpdateUserInfoReq 用户信息更新请求
type UpdateUserInfoReq struct {
	Nickname  string  `json:"nickname"`
	AvatarURL string  `json:"avatar_url"`
	Gender    int32   `json:"gender"`
	Signature string  `json:"signature"`
	Birth     *string `json:"birth"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Ex        string  `json:"ex"`
}

// UpdateUserInfo
func (u *UserService) UpdateUserInfo(ctx context.Context, req UpdateUserInfoReq, userId int64) error {
	var user model.User
	if err := u.db.WithContext(ctx).First(&user, "user_id = ?", userId).Error; err != nil {
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

func (u *UserService) GetUsersPublicInfo(ctx context.Context, userId int64) (dto.UserInfo, error) {
	var user = model.User{}
	if err := u.db.WithContext(ctx).First(&user, "user_id = ?", userId).Error; err != nil {
		return dto.UserInfo{}, err
	}
	return dto.ConvertToUserInfo(user), nil
}

type UserLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *UserService) UserLogin(ctx context.Context, req UserLoginReq) (string, error) {
	var user model.User
	if err := u.db.WithContext(ctx).First(&user, "username = ?", req.Username).Error; err != nil {
		return "", errs.ErrUserNotFound
	}
	//使用加密库比对密码
	bytePassword := []byte(req.Password)
	byteHashedPassword := []byte(user.PasswordHash)
	if bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword) != nil {
		return "", errs.ErrUserPasswordWrong
	}
	token, err := util.GenerateToken(user.UserID)
	if err != nil {
		return "", err
	}
	return token, nil
}
