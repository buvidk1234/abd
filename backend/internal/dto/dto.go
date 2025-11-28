package dto

import (
	"backend/internal/model"
	"time"
)

// UserInfo
type UserInfo struct {
	UserID    string `json:"userID"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarURL"`
	Gender    int32  `json:"gender"`
	Signature string `json:"signature"`
}

func ConvertToUserInfo(user model.User) UserInfo {
	if user.UserID == "" {
		return UserInfo{}
	}
	return UserInfo{
		UserID:    user.UserID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Gender:    user.Gender,
		Signature: user.Signature,
	}
}

// Friend

type FriendRequestInfo struct {
	ID           uint      `json:"id"`
	FromUser     UserInfo  `json:"fromUser"`
	ToUser       UserInfo  `json:"toUser"`
	HandleResult int32     `json:"handleResult"`
	ReqMsg       string    `json:"reqMsg"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	HandleMsg    string    `json:"handleMsg"`
	HandledAt    time.Time `json:"handledAt"`
}

func ConvertToFriendRequestInfo(fr model.FriendRequest, fromUser model.User, toUser model.User) FriendRequestInfo {
	return FriendRequestInfo{
		ID:           fr.ID,
		FromUser:     ConvertToUserInfo(fromUser),
		ToUser:       ConvertToUserInfo(toUser),
		HandleResult: fr.HandleResult,
		ReqMsg:       fr.ReqMsg,
		CreatedAt:    fr.CreatedAt,
		UpdatedAt:    fr.UpdatedAt,
		HandleMsg:    fr.HandleMsg,
		HandledAt:    fr.HandledAt,
	}
}

type FriendInfo struct {
	OwnerUserID string   `json:"ownerUserID"`
	Remark      string   `json:"remark"`
	CreatedAt   int64    `json:"createdAt"`
	FriendUser  UserInfo `json:"friendUser"`
	AddSource   int32    `json:"addSource"`
	IsPinned    bool     `json:"isPinned"`
}

func ConvertToFriendInfo(friend model.Friend, user model.User) FriendInfo {
	return FriendInfo{
		OwnerUserID: friend.OwnerUserID,
		Remark:      friend.Remark,
		CreatedAt:   friend.CreatedAt.Unix(),
		FriendUser:  ConvertToUserInfo(user),
		AddSource:   friend.AddSource,
		IsPinned:    friend.IsPinned,
	}
}
