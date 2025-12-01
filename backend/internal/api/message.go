package api

import (
	"backend/internal/api/apiresp"
	"backend/internal/api/apiresp/errs"
	"backend/internal/service"
	"strconv"

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

func (a *MessageApi) PullSpecifiedMsg(c *gin.Context) {
	userAID, _ := c.Get("my_user_id")
	sessionTypeStr := c.Query("session_type")
	oppositeID := c.Query("opposite_id")
	sessionType, err := strconv.ParseInt(sessionTypeStr, 10, 32)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := a.s.PullSpecifiedMsg(c.Request.Context(), userAID.(string), int32(sessionType), oppositeID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (a *MessageApi) MarkMsgsAsRead(c *gin.Context) {
	var req service.MarkMsgsAsReadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	err := a.s.MarkMsgsAsRead(c.Request.Context(), req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

func (a *MessageApi) PullAllMsg(c *gin.Context) {
	userID, exists := c.Get("my_user_id")
	if !exists {
		apiresp.GinError(c, errs.ErrUnauthorized)
		return
	}
	resp, err := a.s.PullAllMsg(c.Request.Context(), userID.(string))
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
