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

type GroupInfo struct {
	ID                uint      `json:"id"`
	GroupName         string    `json:"groupName"`
	AvatarURL         string    `json:"avatarURL"`
	CreatedAt         time.Time `json:"createdAt"`
	Ex                string    `json:"ex"`
	Status            int32     `json:"status"`
	CreatorUserID     string    `json:"creatorUserID"`
	GroupType         int32     `json:"groupType"`
	NeedVerification  int32     `json:"needVerification"`
	LookMemberInfo    int32     `json:"lookMemberInfo"`
	ApplyMemberFriend int32     `json:"applyMemberFriend"`
}

func ConvertToGroupInfo(group model.Group) GroupInfo {
	return GroupInfo{
		ID:                group.ID,
		GroupName:         group.GroupName,
		AvatarURL:         group.AvatarURL,
		CreatedAt:         group.CreatedAt,
		Ex:                group.Ex,
		Status:            group.Status,
		CreatorUserID:     group.CreatorUserID,
		GroupType:         group.GroupType,
		NeedVerification:  group.NeedVerification,
		LookMemberInfo:    group.LookMemberInfo,
		ApplyMemberFriend: group.ApplyMemberFriend,
	}
}

type GroupMemberInfo struct {
	ID        uint      `json:"id"`
	GroupID   string    `json:"groupID"`
	UserID    string    `json:"userID"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatarURL"`
	RoleLevel int32     `json:"roleLevel"`
	JoinedAt  time.Time `json:"joinedAt"`
}

func ConvertToGroupMemberInfo(member model.GroupMember) GroupMemberInfo {
	return GroupMemberInfo{
		ID:        member.ID,
		GroupID:   member.GroupID,
		UserID:    member.UserID,
		Nickname:  member.Nickname,
		AvatarURL: member.AvatarURL,
		RoleLevel: member.RoleLevel,
		JoinedAt:  member.JoinedAt,
	}
}
