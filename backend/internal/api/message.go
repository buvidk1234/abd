package api

import (
	"backend/internal/api/apiresp"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type MessageApi struct {
	s *service.MessageService
}

func NewMessageApi(s *service.MessageService) *MessageApi {
	return &MessageApi{s: s}
}

func (a *MessageApi) SendMessage(c *gin.Context) {
	var req service.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	err := a.s.SendMessage(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

func (a *MessageApi) PullSpecifiedConv(c *gin.Context) {
	var req service.PullSpecifiedConvReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := a.s.PullSpecifiedConv(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (a *MessageApi) PullConvList(c *gin.Context) {
	var req service.PullConvListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := a.s.PullConvList(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (a *MessageApi) DeleteConversation(c *gin.Context) {
	userID, _ := c.Get("my_user_id")
	conversationID := c.Param("conversation_id")
	err := a.s.DeleteConversation(c.Request.Context(), userID.(string), conversationID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}
