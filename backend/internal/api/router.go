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

	// r.GET("/ws", WsHandler)

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
	userService := service.NewUserService(database.GetDB())
	u := NewUserApi(userService)
	f := NewFriendApi(service.NewFriendService(database.GetDB(), userService))
	g := NewGroupApi(service.NewGroupService(database.GetDB()))
	m := NewMessageApi(service.NewMessageService(database.GetDB()))

	public := r.Group("/user")
	{
		public.POST("/register", u.UserRegister)
		public.POST("/login", u.UserLogin)
	}

	auth := r.Group("/", AuthMiddleware())
	{
		userRouterGroup := auth.Group("/user")
		{
			userRouterGroup.POST("/update-info", u.UpdateUserInfo)
			userRouterGroup.GET("/info", u.GetUsersPublicInfo)
			userRouterGroup.GET("/search", u.SearchUser)
		}

		friendRouterGroup := auth.Group("/friend")
		{
			friendRouterGroup.GET("/search", f.SearchFriend)
			friendRouterGroup.POST("/add", f.ApplyToAddFriend)
			friendRouterGroup.POST("/add-response", f.RespondFriendApply)
			friendRouterGroup.GET("/", f.GetFriendList)
			friendRouterGroup.GET("/black", f.GetPaginationBlacks)
			friendRouterGroup.GET("/:friendId", f.GetSpecifiedFriendsInfo)
			friendRouterGroup.DELETE("/:friendId", f.DeleteFriend)
			friendRouterGroup.POST("/add_black", f.AddBlack)
			friendRouterGroup.POST("/remove_black", f.RemoveBlack)
			friendRouterGroup.GET("/apply", f.GetFriendApplyList)
			friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList)
		}

		groupRouterGroup := auth.Group("/groups")
		{
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
		msgGroup := auth.Group("/msg")
		{
			msgGroup.POST("/send", m.SendMessage)              // 发送消息
			msgGroup.GET("/pull", m.PullConvList)              // 拉取会话列表
			msgGroup.GET("/pull/:convID", m.PullSpecifiedConv) // 拉取某个会话的消息
			// msgGroup.POST("/newest-seq", m.GetSeq)
			// msgGroup.POST("/search", m.SearchMsg)
			// msgGroup.POST("/send", m.SendMessage)
			// msgGroup.POST("/send-business-notification", m.SendBusinessNotification)
			// msgGroup.POST("/pull", m.PullMsgBySeqs)
			// msgGroup.POST("/revoke", m.RevokeMsg)
			// msgGroup.POST("/mark-read", m.MarkMsgsAsRead)
			// msgGroup.POST("/sync-convs", m.GetConversationsHasReadAndMaxSeq)
			// msgGroup.POST("/read", m.SetConversationHasReadSeq)

			// msgGroup.POST("/clear-conv", m.ClearConversationsMsg)
			// msgGroup.POST("/clear-all", m.UserClearAllMsg)
			// msgGroup.POST("/delete", m.DeleteMsgs)
			// msgGroup.POST("/delete-physical", m.DeleteMsgPhysical)

			// msgGroup.POST("/batch-send", m.BatchSendMsg)
			// msgGroup.POST("/server-time", m.GetServerTime)
		}

	}
	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
