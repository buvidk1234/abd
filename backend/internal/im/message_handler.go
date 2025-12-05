package im

import (
	"backend/internal/service"
	"context"
	"encoding/json"
	"log"
	"sync"
)

const (
	TextPing = "ping"
	TextPong = "pong"
)

type TextMessage struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type Req struct {
	ReqIdentifier int32           `json:"req_identifier,string" validate:"required"`
	Token         string          `json:"token"`
	SendID        string          `json:"send_id"`
	MsgIncr       string          `json:"msg_incr"`
	Data          json.RawMessage `json:"data"`
}

var reqPool = sync.Pool{
	New: func() any {
		return new(Req)
	},
}

func getReq() *Req {
	req := reqPool.Get().(*Req)
	req.Data = nil
	req.MsgIncr = ""
	req.ReqIdentifier = 0
	req.SendID = ""
	req.Token = ""
	return req
}
func freeReq(req *Req) {
	reqPool.Put(req)
}

type Resp struct {
	ReqIdentifier int32  `json:"req_identifier"`
	MsgIncr       string `json:"msg_incr"`
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	Data          any    `json:"data"`
}

type MessageHandler interface {
	SendMessage(ctx context.Context, data *Req) (any, error)
	PullSpecifiedConv(ctx context.Context, data *Req) (any, error)
	PullConvList(ctx context.Context, data *Req) (any, error)
	// UserLogout(ctx context.Context, data *Req) ([]byte, error)
}

var _ MessageHandler = (*ServiceHandler)(nil)

type ServiceHandler struct {
	messageService *service.MessageService
	// pushClient *
}

func NewServiceHandler(messageService *service.MessageService) *ServiceHandler {
	return &ServiceHandler{
		messageService: messageService,
	}
}

func (s *ServiceHandler) SendMessage(ctx context.Context, data *Req) (any, error) {
	// encode
	log.Print("SendMessage")
	var sendMsgReq service.SendMessageReq
	if err := json.Unmarshal(data.Data, &sendMsgReq); err != nil {
		return nil, err
	}
	err := s.messageService.SendMessage(ctx, sendMsgReq)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ServiceHandler) PullSpecifiedConv(ctx context.Context, data *Req) (any, error) {
	var pullReq service.PullSpecifiedConvReq
	if err := json.Unmarshal(data.Data, &pullReq); err != nil {
		return nil, err
	}
	resp, err := s.messageService.PullSpecifiedConv(ctx, pullReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (s *ServiceHandler) PullConvList(ctx context.Context, data *Req) (any, error) {
	var pullReq service.PullConvListReq
	if err := json.Unmarshal(data.Data, &pullReq); err != nil {
		return nil, err
	}
	resp, err := s.messageService.PullConvList(ctx, pullReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
