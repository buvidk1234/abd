package api

import (
	"backend/internal/api/apiresp"
	"backend/internal/api/apiresp/errs"
	"backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendApi struct {
	friendService *service.FriendService
}

func NewFriendApi(friendService *service.FriendService) *FriendApi {
	return &FriendApi{friendService: friendService}
}

// 申请添加好友
func (f *FriendApi) ApplyToAddFriend(c *gin.Context) {
	var req service.ApplyToAddFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.FromUserID = c.GetInt64("user_id")
	err := f.friendService.ApplyToAddFriend(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// 响应好友申请
func (f *FriendApi) RespondFriendApply(c *gin.Context) {
	var req service.RespondFriendApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.HandlerUserID = c.GetInt64("user_id")
	err := f.friendService.RespondFriendApply(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// 获取好友列表
func (f *FriendApi) GetFriendList(c *gin.Context) {
	var req service.GetPaginationFriendsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	resp, err := f.friendService.GetPaginationFriends(c.Request.Context(), req, c.GetInt64("user_id"))
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// 获取指定好友信息
func (f *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	friendUserIDStr := c.Param("friendId")
	friendUserID, err := strconv.ParseInt(friendUserIDStr, 10, 64)
	if err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	friendInfo, err := f.friendService.GetSpecifiedFriendInfo(c.Request.Context(), c.GetInt64("user_id"), friendUserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, friendInfo)
}

// 删除好友
func (f *FriendApi) DeleteFriend(c *gin.Context) {
	friendUserIDStr := c.Param("friendId")
	friendUserID, err := strconv.ParseInt(friendUserIDStr, 10, 64)
	if err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	err = f.friendService.DeleteFriend(c.Request.Context(), c.GetInt64("user_id"), friendUserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// 获取收到的好友申请列表
func (f *FriendApi) GetFriendApplyList(c *gin.Context) {
	var req service.GetPaginationFriendApplyListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.ToUserID = c.GetInt64("user_id")
	resp, err := f.friendService.GetPaginationFriendApplyList(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// 获取自己发出的好友申请列表
func (f *FriendApi) GetSelfApplyList(c *gin.Context) {
	var req service.GetPaginationSelfFriendApplyListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.FromUserID = c.GetInt64("user_id")
	resp, err := f.friendService.GetPaginationSelfApplyList(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// 添加黑名单
func (f *FriendApi) AddBlack(c *gin.Context) {
	var req service.AddBlackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.OwnerUserID = c.GetInt64("user_id")
	req.OperatorUserID = req.OwnerUserID
	if err := f.friendService.AddBlack(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// 移除黑名单
func (f *FriendApi) RemoveBlack(c *gin.Context) {
	var req service.RemoveBlackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.OwnerUserID = c.GetInt64("user_id")
	if err := f.friendService.RemoveBlack(c.Request.Context(), req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// 获取黑名单列表
func (f *FriendApi) GetPaginationBlacks(c *gin.Context) {
	var req service.GetPaginationBlacksReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiresp.GinError(c, errs.ErrInvalidParam)
		return
	}
	req.OwnerUserID = c.GetInt64("user_id")
	resp, err := f.friendService.GetPaginationBlacks(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}
