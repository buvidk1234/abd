package api

import (
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FriendApi struct {
	friendService *service.FriendService
}

func NewFriendApi(friendService service.FriendService) *FriendApi {
	return &FriendApi{friendService: &friendService}
}

// 申请添加好友
func (f *FriendApi) ApplyToAddFriend(c *gin.Context) {
	var req service.ApplyToAddFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
	}
	err := f.friendService.ApplyToAddFriend(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"msg": "好友申请已发送"})
}

// 响应好友申请
func (f *FriendApi) RespondFriendApply(c *gin.Context) {
	var req service.RespondFriendApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
	}
	err := f.friendService.RespondFriendApply(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"msg": "好友申请已处理"})
}

// 获取好友列表
func (f *FriendApi) GetFriendList(c *gin.Context) {
	var req service.GetPaginationFriendsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	resp, err := f.friendService.GetPaginationFriends(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resp})
}

// 获取指定好友信息
func (f *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	ownerUserID := c.PostForm("ownerUserID")
	friendUserID := c.PostForm("friendUserID")
	friendInfo, err := f.friendService.GetSpecifiedFriendInfo(c.Request.Context(), ownerUserID, friendUserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"friendInfo": friendInfo})
}

// 删除好友
func (f *FriendApi) DeleteFriend(c *gin.Context) {
	var req service.DeleteFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	err := f.friendService.DeleteFriend(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "好友已删除"})
}

// 获取收到的好友申请列表
func (f *FriendApi) GetFriendApplyList(c *gin.Context) {
	var req service.GetPaginationFriendApplyListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	resp, err := f.friendService.GetPaginationFriendApplyList(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resp})
}

// 获取自己发出的好友申请列表
func (f *FriendApi) GetSelfApplyList(c *gin.Context) {
	var req service.GetPaginationSelfFriendApplyListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	resp, err := f.friendService.GetPaginationSelfApplyList(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resp})
}

// 添加黑名单
func (f *FriendApi) AddBlack(c *gin.Context) {

}

// 移除黑名单
func (f *FriendApi) RemoveBlack(c *gin.Context) {

}

// 获取黑名单列表
func (f *FriendApi) GetPaginationBlacks(c *gin.Context) {

}
