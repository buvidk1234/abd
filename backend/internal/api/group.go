package api

import (
	"backend/internal/api/apiresp"
	"backend/internal/api/apiresp/errs"
	"backend/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type GroupApi struct {
	s *service.GroupService
}

func NewGroupApi(s *service.GroupService) *GroupApi {
	return &GroupApi{s: s}
}

func (a *GroupApi) CreateGroup(c *gin.Context) {
	var req service.CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	groupID, err := a.s.CreateGroup(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"groupID": groupID})
}

func (a *GroupApi) GetGroupsInfo(c *gin.Context) {
	idsStr := c.Query("ids")
	var ids []string
	if idsStr != "" {
		ids = strings.Split(idsStr, ",")
	}
	groupInfos, err := a.s.GetGroupsInfo(c.Request.Context(), ids)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"groupInfos": groupInfos})
}

func (a *GroupApi) GetGroupMemberList(c *gin.Context) {
	id := c.Param("id")
	memberList, err := a.s.GetGroupMemberList(c.Request.Context(), id)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"memberList": memberList})
}

func (a *GroupApi) JoinGroup(c *gin.Context) {
	var req service.JoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	joined, err := a.s.JoinGroup(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	status := "pending"
	if joined {
		status = "joined"
	}
	apiresp.GinSuccess(c, gin.H{"status": status})
}

func (a *GroupApi) QuitGroup(c *gin.Context) {
	var req service.QuitGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	req.UserID = c.Param("userID")
	if err := a.s.QuitGroup(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"msg": "已退出群组"})
}

func (a *GroupApi) InviteUserToGroup(c *gin.Context) {
	var req service.InviteUserToGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	joined, err := a.s.InviteUserToGroup(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	status := "pending"
	if joined {
		status = "joined"
	}
	apiresp.GinSuccess(c, gin.H{"status": status})
}

func (a *GroupApi) KickGroupMember(c *gin.Context) {
	var req service.KickGroupMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	req.TargetUserID = c.Param("userID")
	if err := a.s.KickGroupMember(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"msg": "已移除群成员"})
}

func (a *GroupApi) DismissGroup(c *gin.Context) {
	var req service.DismissGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	if err := a.s.DismissGroup(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"msg": "群组已解散"})
}

func (a *GroupApi) SetGroupInfo(c *gin.Context) {
	var req service.SetGroupInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	if err := a.s.SetGroupInfo(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"msg": "群信息已更新"})
}

func (a *GroupApi) SetGroupMemberInfo(c *gin.Context) {
	var req service.SetGroupMemberInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.GroupID = c.Param("id")
	req.UserID = c.Param("userID")
	if err := a.s.SetGroupMemberInfo(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, gin.H{"msg": "群成员信息已更新"})
}
