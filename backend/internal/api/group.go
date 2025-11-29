package api

import (
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type GroupApi struct {
	s *service.GroupService
}

func NewGroupApi(s *service.GroupService) *GroupApi {
	return &GroupApi{s: s}
}

func (a *GroupApi) CreateGroup(c *gin.Context) {

}

func (a *GroupApi) GetGroupsInfo(c *gin.Context) {

}

func (a *GroupApi) GetGroupMemberList(c *gin.Context) {

}

func (a *GroupApi) JoinGroup(c *gin.Context) {

}

func (a *GroupApi) QuitGroup(c *gin.Context) {

}

func (a *GroupApi) InviteUserToGroup(c *gin.Context) {

}

func (a *GroupApi) KickGroupMember(c *gin.Context) {

}

func (a *GroupApi) DismissGroup(c *gin.Context) {

}

func (a *GroupApi) SetGroupInfo(c *gin.Context) {

}

func (a *GroupApi) SetGroupMemberInfo(c *gin.Context) {

}
