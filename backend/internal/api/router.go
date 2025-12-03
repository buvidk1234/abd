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
	u := NewUserApi(service.NewUserService(database.GetDB()))
	f := NewFriendApi(service.NewFriendService(database.GetDB()))
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
		}

		friendRouterGroup := auth.Group("/friend")
		{
			friendRouterGroup.POST("/add", f.ApplyToAddFriend)
			friendRouterGroup.POST("/add-response", f.RespondFriendApply)
			friendRouterGroup.GET("/", f.GetFriendList)
			friendRouterGroup.GET("/black", f.GetPaginationBlacks)
			friendRouterGroup.GET("/:friendId", f.GetSpecifiedFriendsInfo)
			friendRouterGroup.DELETE("/:friendId", f.DeleteFriend)
			friendRouterGroup.POST("/add_black", f.AddBlack)
			friendRouterGroup.POST("/remove_black", f.RemoveBlack)
			friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)
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
			// msgGroup.DELETE("/delete_conversation", m.DeleteConversation) // 删除会话
			// msgGroup.POST("/revoke_msg", m.RevokeMsg)             // 撤回消息
			// msgGroup.POST("/delete_msgs", m.DeleteMsgs)           // 删除消息
			// msgGroup.POST("/newest_seq", m.GetSeq)                // 获取最新消息序列号
		}

	}
	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
