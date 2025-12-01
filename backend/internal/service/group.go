package service

import (
	"backend/internal/api/apiresp/errs"
	"backend/internal/dto"
	"backend/internal/model"
	"context"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type GroupService struct {
	db *gorm.DB
}

const (
	roleOwner  int32 = 100
	roleAdmin  int32 = 60
	roleMember int32 = 10
)

type CreateGroupReq struct {
	GroupName         string `json:"groupName" binding:"required"`
	AvatarURL         string `json:"avatarURL" binding:"required"`
	Ex                string `json:"ex"`
	CreatorUserID     string `json:"creatorUserID" binding:"required"`
	GroupType         int32  `json:"groupType" default:"1"`
	NeedVerification  int32  `json:"needVerification" default:"1"`
	LookMemberInfo    int32  `json:"lookMemberInfo" default:"1"`
	ApplyMemberFriend int32  `json:"applyMemberFriend" default:"1"`
}

type JoinGroupReq struct {
	GroupID       string `json:"groupID" binding:"required"`
	UserID        string `json:"userID" binding:"required"`
	ReqMsg        string `json:"reqMsg"`
	JoinSource    int32  `json:"joinSource"`
	InviterUserID string `json:"inviterUserID"`
}

type QuitGroupReq struct {
	GroupID        string `json:"groupID" binding:"required"`
	UserID         string `json:"userID" binding:"required"`
	OperatorUserID string `json:"operatorUserID" binding:"required"`
}

type InviteUserToGroupReq struct {
	GroupID       string `json:"groupID" binding:"required"`
	InviterUserID string `json:"inviterUserID" binding:"required"`
	InviteeUserID string `json:"inviteeUserID" binding:"required"`
}

type KickGroupMemberReq struct {
	GroupID        string `json:"groupID" binding:"required"`
	OperatorUserID string `json:"operatorUserID" binding:"required"`
	TargetUserID   string `json:"targetUserID" binding:"required"`
}

type DismissGroupReq struct {
	GroupID        string `json:"groupID" binding:"required"`
	OperatorUserID string `json:"operatorUserID" binding:"required"`
}

type SetGroupInfoReq struct {
	GroupID           string  `json:"groupID" binding:"required"`
	OperatorUserID    string  `json:"operatorUserID" binding:"required"`
	GroupName         *string `json:"groupName"`
	AvatarURL         *string `json:"avatarURL"`
	Notification      *string `json:"notification"`
	Introduction      *string `json:"introduction"`
	NeedVerification  *int32  `json:"needVerification"`
	LookMemberInfo    *int32  `json:"lookMemberInfo"`
	ApplyMemberFriend *int32  `json:"applyMemberFriend"`
}

type SetGroupMemberInfoReq struct {
	GroupID        string  `json:"groupID" binding:"required"`
	UserID         string  `json:"userID" binding:"required"`
	OperatorUserID string  `json:"operatorUserID" binding:"required"`
	Nickname       *string `json:"nickname"`
	AvatarURL      *string `json:"avatarURL"`
	RoleLevel      *int32  `json:"roleLevel"`
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) CreateGroup(ctx context.Context, req CreateGroupReq) (string, error) {
	var groupID string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		group := model.Group{
			GroupName:         req.GroupName,
			AvatarURL:         req.AvatarURL,
			Ex:                req.Ex,
			Status:            1, // 1正常 2解散
			CreatorUserID:     req.CreatorUserID,
			GroupType:         req.GroupType,
			NeedVerification:  req.NeedVerification,
			LookMemberInfo:    req.LookMemberInfo,
			ApplyMemberFriend: req.ApplyMemberFriend,
		}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		groupID = strconv.Itoa(int(group.ID))
		creatorMember := model.GroupMember{
			GroupID:        groupID,
			UserID:         req.CreatorUserID,
			RoleLevel:      roleOwner,
			JoinSource:     1,
			OperatorUserID: req.CreatorUserID,
		}
		if err := tx.Create(&creatorMember).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return groupID, nil
}

func (s *GroupService) GetGroupsInfo(ctx context.Context, ids []string) ([]dto.GroupInfo, error) {
	var groups []model.Group
	db := s.db.WithContext(ctx).Model(&model.Group{})
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}
	if err := db.Find(&groups).Error; err != nil {
		return nil, err
	}
	groupInfos := make([]dto.GroupInfo, 0, len(groups))
	for _, group := range groups {
		groupInfos = append(groupInfos, dto.ConvertToGroupInfo(group))
	}
	return groupInfos, nil
}

func (s *GroupService) GetGroupMemberList(ctx context.Context, id string) ([]dto.GroupMemberInfo, error) {
	var members []model.GroupMember
	if err := s.db.WithContext(ctx).Where("group_id = ?", id).Find(&members).Error; err != nil {
		return nil, err
	}
	memberInfos := make([]dto.GroupMemberInfo, 0, len(members))
	for _, member := range members {
		memberInfos = append(memberInfos, dto.ConvertToGroupMemberInfo(member))
	}
	return memberInfos, nil
}

func (s *GroupService) JoinGroup(ctx context.Context, req JoinGroupReq) (bool, error) {
	group, err := s.getGroup(ctx, req.GroupID)
	if err != nil {
		return false, err
	}
	var memberCount int64
	if err := s.db.WithContext(ctx).Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", req.GroupID, req.UserID).
		Count(&memberCount).Error; err != nil {
		return false, err
	}
	if memberCount > 0 {
		return true, nil
	}
	if group.NeedVerification == 1 {
		var pending int64
		if err := s.db.WithContext(ctx).Model(&model.GroupRequest{}).
			Where("group_id = ? AND user_id = ? AND handle_result = 0", req.GroupID, req.UserID).
			Count(&pending).Error; err != nil {
			return false, err
		}
		if pending == 0 {
			request := model.GroupRequest{
				UserID:        req.UserID,
				GroupID:       req.GroupID,
				ReqMsg:        req.ReqMsg,
				JoinSource:    req.JoinSource,
				InviterUserID: req.InviterUserID,
				HandleResult:  0,
			}
			if err := s.db.WithContext(ctx).Create(&request).Error; err != nil {
				return false, err
			}
		}
		return false, nil
	}
	member := model.GroupMember{
		GroupID:        req.GroupID,
		UserID:         req.UserID,
		RoleLevel:      roleMember,
		JoinSource:     req.JoinSource,
		InviterUserID:  req.InviterUserID,
		OperatorUserID: req.UserID,
	}
	if err := s.db.WithContext(ctx).Create(&member).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (s *GroupService) QuitGroup(ctx context.Context, req QuitGroupReq) error {
	if req.OperatorUserID != req.UserID {
		return errs.ErrGroupQuitSelfOnly
	}
	group, err := s.getGroup(ctx, req.GroupID)
	if err != nil {
		return err
	}
	if group.CreatorUserID == req.UserID {
		return errs.ErrGroupOwnerCannotQuit
	}
	result := s.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", req.GroupID, req.UserID).Delete(&model.GroupMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrGroupMemberNotFound
	}
	return nil
}

func (s *GroupService) InviteUserToGroup(ctx context.Context, req InviteUserToGroupReq) (bool, error) {
	if _, err := s.getGroup(ctx, req.GroupID); err != nil {
		return false, err
	}
	if _, err := s.getMember(ctx, req.GroupID, req.InviterUserID); err != nil {
		return false, err
	}
	return s.JoinGroup(ctx, JoinGroupReq{
		GroupID:       req.GroupID,
		UserID:        req.InviteeUserID,
		InviterUserID: req.InviterUserID,
		JoinSource:    2,
	})
}

func (s *GroupService) KickGroupMember(ctx context.Context, req KickGroupMemberReq) error {
	group, err := s.getGroup(ctx, req.GroupID)
	if err != nil {
		return err
	}
	operator, err := s.getMember(ctx, req.GroupID, req.OperatorUserID)
	if err != nil {
		return err
	}
	target, err := s.getMember(ctx, req.GroupID, req.TargetUserID)
	if err != nil {
		return err
	}
	if req.TargetUserID == group.CreatorUserID || target.RoleLevel == roleOwner {
		return errs.ErrGroupCannotRemoveOwner
	}
	if operator.UserID == target.UserID {
		return errs.ErrGroupUseQuitGroup
	}
	if operator.RoleLevel < roleAdmin {
		return errs.ErrGroupPermissionDenied.WithDetail("没有足够的权限踢人")
	}
	if operator.RoleLevel <= target.RoleLevel {
		return errs.ErrGroupPermissionDenied.WithDetail("权限不足，无法移除该成员")
	}
	return s.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", req.GroupID, req.TargetUserID).Delete(&model.GroupMember{}).Error
}

func (s *GroupService) DismissGroup(ctx context.Context, req DismissGroupReq) error {
	group, err := s.getGroup(ctx, req.GroupID)
	if err != nil {
		return err
	}
	if group.CreatorUserID != req.OperatorUserID {
		return errs.ErrGroupPermissionDenied.WithDetail("只有群主可以解散群组")
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Group{}).Where("id = ?", group.ID).Update("status", 2).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", req.GroupID).Delete(&model.GroupMember{}).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", req.GroupID).Delete(&model.GroupRequest{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *GroupService) SetGroupInfo(ctx context.Context, req SetGroupInfoReq) error {
	group, err := s.getGroup(ctx, req.GroupID)
	if err != nil {
		return err
	}
	operator, err := s.getMember(ctx, req.GroupID, req.OperatorUserID)
	if err != nil {
		return err
	}
	if operator.RoleLevel < roleAdmin && req.OperatorUserID != group.CreatorUserID {
		return errs.ErrGroupPermissionDenied.WithDetail("没有权限修改群信息")
	}
	updates := map[string]interface{}{}
	if req.GroupName != nil {
		updates["group_name"] = *req.GroupName
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}
	if req.Notification != nil {
		updates["notification"] = *req.Notification
	}
	if req.Introduction != nil {
		updates["introduction"] = *req.Introduction
	}
	if req.NeedVerification != nil {
		updates["need_verification"] = *req.NeedVerification
	}
	if req.LookMemberInfo != nil {
		updates["look_member_info"] = *req.LookMemberInfo
	}
	if req.ApplyMemberFriend != nil {
		updates["apply_member_friend"] = *req.ApplyMemberFriend
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Model(&model.Group{}).Where("id = ?", group.ID).Updates(updates).Error
}

func (s *GroupService) SetGroupMemberInfo(ctx context.Context, req SetGroupMemberInfoReq) error {
	if _, err := s.getGroup(ctx, req.GroupID); err != nil {
		return err
	}
	operator, err := s.getMember(ctx, req.GroupID, req.OperatorUserID)
	if err != nil {
		return err
	}
	member, err := s.getMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return err
	}
	if req.OperatorUserID != req.UserID && operator.RoleLevel < roleAdmin {
		return errs.ErrGroupPermissionDenied.WithDetail("没有权限修改其他成员信息")
	}
	updates := map[string]interface{}{}
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}
	if req.RoleLevel != nil {
		if operator.RoleLevel < roleOwner {
			return errs.ErrGroupOnlyOwnerCanSetRole
		}
		if *req.RoleLevel >= operator.RoleLevel {
			return errs.ErrGroupRoleLevelTooHigh
		}
		updates["role_level"] = *req.RoleLevel
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Model(&member).Updates(updates).Error
}

func (s *GroupService) getGroup(ctx context.Context, groupID string) (model.Group, error) {
	var group model.Group
	if err := s.db.WithContext(ctx).First(&group, groupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Group{}, errs.ErrGroupNotFound
		}
		return model.Group{}, err
	}
	if group.Status == 2 {
		return model.Group{}, errs.ErrGroupDismissed
	}
	return group, nil
}

func (s *GroupService) getMember(ctx context.Context, groupID, userID string) (model.GroupMember, error) {
	var member model.GroupMember
	if err := s.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.GroupMember{}, errs.ErrGroupMemberNotFound
		}
		return model.GroupMember{}, err
	}
	return member, nil
}
