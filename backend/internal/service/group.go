package service

import (
	"context"

	"gorm.io/gorm"
)

type GroupService struct {
	db *gorm.DB
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) CreateGroup(ctx context.Context) {

}

func (s *GroupService) GetGroupsInfo(ctx context.Context) {

}

func (s *GroupService) GetGroupMemberList(ctx context.Context) {

}

func (s *GroupService) JoinGroup(ctx context.Context) {

}

func (s *GroupService) QuitGroup(ctx context.Context) {

}

func (s *GroupService) InviteUserToGroup(ctx context.Context) {

}

func (s *GroupService) KickGroupMember(ctx context.Context) {

}

func (s *GroupService) DismissGroup(ctx context.Context) {

}

func (s *GroupService) SetGroupInfo(ctx context.Context) {

}

func (s *GroupService) SetGroupMemberInfo(ctx context.Context) {

}
