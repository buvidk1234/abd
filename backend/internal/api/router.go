package api

import (
	"backend/internal/service"

	"backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ws", WsHandler)

	// middlewares
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// user routes
	u := NewUserApi(service.NewUserService(database.GetDB()))
	{
		userRouterGroup := r.Group("/user")
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", u.UpdateUserInfo)
		userRouterGroup.POST("/get_users_info", u.GetUsersPublicInfo)
		userRouterGroup.POST("/user_login", u.UserLogin)
	}
	// friend routes
	f := NewFriendApi(service.NewFriendService(database.GetDB()))
	{
		friendRouterGroup := r.Group("/friend")
		friendRouterGroup.POST("/add_friend", f.ApplyToAddFriend)
		friendRouterGroup.POST("/add_friend_response", f.RespondFriendApply)
		friendRouterGroup.POST("/get_friend_list", f.GetFriendList)
		friendRouterGroup.POST("/get_specified_friends_info", f.GetSpecifiedFriendsInfo)
		friendRouterGroup.POST("/delete_friend", f.DeleteFriend)
		friendRouterGroup.POST("/add_black", f.AddBlack)
		friendRouterGroup.POST("/remove_black", f.RemoveBlack)
		friendRouterGroup.POST("/get_black_list", f.GetPaginationBlacks)
		friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)
		friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList)
	}

	// Group
	g := NewGroupApi(service.NewGroupService(database.GetDB()))
	{
		groupRouterGroup := r.Group("/group")
		groupRouterGroup.POST("/create_group", g.CreateGroup)                 // 创建群组
		groupRouterGroup.POST("/get_groups_info", g.GetGroupsInfo)            // 获取群组信息
		groupRouterGroup.POST("/get_group_member_list", g.GetGroupMemberList) // 获取群成员列表
		groupRouterGroup.POST("/join_group", g.JoinGroup)                     // 加入群组
		groupRouterGroup.POST("/quit_group", g.QuitGroup)                     // 退出群组
		groupRouterGroup.POST("/invite_user_to_group", g.InviteUserToGroup)   // 邀请进群
		groupRouterGroup.POST("/kick_group", g.KickGroupMember)               // 踢人
		groupRouterGroup.POST("/dismiss_group", g.DismissGroup)               // 解散群组
		groupRouterGroup.POST("/set_group_info", g.SetGroupInfo)              // 设置群信息
		groupRouterGroup.POST("/set_group_member_info", g.SetGroupMemberInfo) // 设置群成员信息
	}

	return r
}
