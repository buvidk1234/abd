package api

import (
	"backend/docs"
	"backend/internal/service"
	"time"

	"backend/internal/pkg/database"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ws", WsHandler)

	// swag init -g cmd/server/main.go -d ./ -o ./docs
	docs.SwaggerInfo.BasePath = "/"

	// middlewares
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))
	// auth
	r.Use(AuthMiddleware())

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
		groupRouterGroup := r.Group("/groups")
		groupRouterGroup.POST("", g.CreateGroup)                                // 创建群组
		groupRouterGroup.GET("", g.GetGroupsInfo)                               // 获取群组信息
		groupRouterGroup.GET("/:id/members", g.GetGroupMemberList)              // 获取群成员列表
		groupRouterGroup.POST("/:id/join", g.JoinGroup)                         // 加入群组/申请进群
		groupRouterGroup.DELETE("/:id/members/:userID", g.QuitGroup)            // 退出群组
		groupRouterGroup.POST("/:id/invitations", g.InviteUserToGroup)          // 邀请进群
		groupRouterGroup.DELETE("/:id/members/:userID/kick", g.KickGroupMember) // 踢人
		groupRouterGroup.DELETE("/:id", g.DismissGroup)                         // 解散群组
		groupRouterGroup.POST("/:id", g.SetGroupInfo)                           // 设置群信息
		groupRouterGroup.POST("/:id/members/:userID", g.SetGroupMemberInfo)     // 设置群成员信息
	}

	// Message
	m := NewMessageApi(service.NewMessageService(database.GetDB()))
	{
		msgGroup := r.Group("/msg")
		msgGroup.POST("/send_msg", m.SendMessage)                  // 发送消息
		msgGroup.POST("/pull_specified_conv", m.PullSpecifiedConv) // 拉取某个会话的消息
		msgGroup.POST("/pull_conv_list", m.PullConvList)           // 拉取会话列表
		// msgGroup.DELETE("/delete_conversation", m.DeleteConversation) // 删除会话
		// msgGroup.POST("/revoke_msg", m.RevokeMsg)             // 撤回消息
		// msgGroup.POST("/delete_msgs", m.DeleteMsgs)           // 删除消息
		// msgGroup.POST("/newest_seq", m.GetSeq)                // 获取最新消息序列号
	}

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
