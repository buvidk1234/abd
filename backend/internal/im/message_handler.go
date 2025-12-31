package im

import (
	"backend/internal/pkg/kafka"
	"backend/internal/pkg/prommetrics"
	"backend/internal/service"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

const (
	TextPing = "ping"
	TextPong = "pong"
)

type TextMessage struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type InboundReq struct {
	ReqIdentifier int32           `json:"req_identifier" validate:"required"`
	MsgIncr       string          `json:"msg_incr"`
	Data          json.RawMessage `json:"data"`
}

type Req struct {
	InboundReq
	Token  string
	SendID int64
}

var reqPool = sync.Pool{
	New: func() any {
		return new(Req)
	},
}

func getReq(token string, sendId int64) *Req {
	req := reqPool.Get().(*Req)
	req.Data = nil
	req.MsgIncr = ""
	req.ReqIdentifier = 0
	req.SendID = sendId
	req.Token = token
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

//	type MessageHandler interface {
//		SendMessage(ctx context.Context, data *Req) (any, error)
//		PullSpecifiedConv(ctx context.Context, data *Req) (any, error)
//		PullConvList(ctx context.Context, data *Req) (any, error)
//		// UserLogout(ctx context.Context, data *Req) ([]byte, error)
//	}
type MessageHandler interface {
	GetSeq(ctx context.Context, data *Req) (any, error)
	SendMessage(ctx context.Context, data *Req) (any, error)
	PullMessageBySeqList(ctx context.Context, data *Req) (any, error)
	GetConversationsHasReadAndMaxSeq(ctx context.Context, data *Req) (any, error)
	GetSeqMessage(ctx context.Context, data *Req) (any, error)
	GetLastMessage(ctx context.Context, data *Req) (any, error)
}

var _ MessageHandler = (*ServiceHandler)(nil)

type ServiceHandler struct {
	messageService *service.MessageService
	producer       sarama.SyncProducer
	// pushClient *
}

func NewServiceHandler(messageService *service.MessageService, producer sarama.SyncProducer) *ServiceHandler {
	return &ServiceHandler{
		messageService: messageService,
		producer:       producer,
	}
}

func (s *ServiceHandler) GetSeq(ctx context.Context, data *Req) (any, error) {
	log.Printf("GetSeq request: %d", data.SendID)
	resp, err := s.messageService.GetMaxSeq(ctx, data.SendID)
	if err != nil {
		return nil, err
	}
	log.Printf("GetSeq response: %+v", resp)
	// r, _ := json.Marshal(resp)
	return resp, nil
}

func (s *ServiceHandler) SendMessage(ctx context.Context, data *Req) (any, error) {
	// encode
	log.Printf("SendMessage: %+v", data)
	// var sendMsgReq service.SendMessageReq
	// if err := json.Unmarshal(data.Data, &sendMsgReq); err != nil {
	// 	return nil, err
	// }
	// err := s.messageService.SendMessage(ctx, sendMsgReq)
	// if err != nil {
	// 	return nil, err
	// }
	// return nil, nil
	msg := &sarama.ProducerMessage{
		Topic: kafka.ComingMessageTopic,
		Value: sarama.ByteEncoder(data.Data),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		prommetrics.MsgProcessFailedCounter.Inc()
		log.Printf("FAILED to send kafka message: %v", err)
		return nil, err
	}
	prommetrics.MsgProcessSuccessCounter.Inc()
	log.Printf("message sent to partition=%d offset=%d", partition, offset)
	return nil, nil
}
func (s *ServiceHandler) PullMessageBySeqList(ctx context.Context, data *Req) (any, error) {
	var pullReq service.PullMessageBySeqsReq
	if err := json.Unmarshal(data.Data, &pullReq); err != nil {
		return nil, err
	}
	log.Printf("PullMessageBySeqList request: %+v", pullReq)
	resp, err := s.messageService.PullMessageBySeqs(ctx, data.SendID, pullReq)
	if err != nil {
		return nil, err
	}
	log.Printf("PullMessageBySeqList response: %+v", resp)
	// r, _ := json.Marshal(resp)
	return resp, nil
}

func (s *ServiceHandler) GetConversationsHasReadAndMaxSeq(ctx context.Context, data *Req) (any, error) {
	var getReq service.GetConversationsHasReadAndMaxSeqReq
	if err := json.Unmarshal(data.Data, &getReq); err != nil {
		return nil, err
	}
	log.Printf("GetConversationsHasReadAndMaxSeq request: %+v", getReq)
	resp, err := s.messageService.GetConversationsHasReadAndMaxSeq(ctx, getReq)
	if err != nil {
		return nil, err
	}
	log.Printf("GetConversationsHasReadAndMaxSeq response: %+v", resp)
	// r, _ := json.Marshal(resp)
	return resp, nil
}
func (s *ServiceHandler) GetLastMessage(ctx context.Context, data *Req) (any, error) {
	var getLastMsgReq service.GetLastMessageReq
	if err := json.Unmarshal(data.Data, &getLastMsgReq); err != nil {
		return nil, err
	}
	log.Printf("GetLastMessage request: %+v", getLastMsgReq)
	resp, err := s.messageService.GetLastMessage(ctx, getLastMsgReq)
	if err != nil {
		return nil, err
	}
	log.Printf("GetLastMessage response: %+v", resp)
	// r, _ := json.Marshal(resp)
	return resp, nil
}

func (s *ServiceHandler) GetSeqMessage(ctx context.Context, data *Req) (any, error) {
	var getSeqMsgReq service.GetSeqMessageReq
	if err := json.Unmarshal(data.Data, &getSeqMsgReq); err != nil {
		return nil, err
	}
	log.Printf("GetSeqMessage request: %+v", getSeqMsgReq)
	resp, err := s.messageService.GetSeqMessage(ctx, getSeqMsgReq)
	if err != nil {
		return nil, err
	}
	log.Printf("GetSeqMessage response: %+v", resp)
	// r, _ := json.Marshal(resp)
	return resp, nil
}

func (s *ServiceHandler) PullSpecifiedConv(ctx context.Context, data *Req) (any, error) {
	var pullReq service.PullSpecifiedConvReq
	if err := json.Unmarshal(data.Data, &pullReq); err != nil {
		return nil, err
	}
	log.Printf("PullSpecifiedConv request: %+v", pullReq)
	resp, err := s.messageService.PullSpecifiedConv(ctx, pullReq)
	if err != nil {
		return nil, err
	}
	log.Printf("PullSpecifiedConv response: %+v", resp)
	return resp, nil
}
func (s *ServiceHandler) PullConvList(ctx context.Context, data *Req) (any, error) {
	var pullReq service.PullConvListReq
	if err := json.Unmarshal(data.Data, &pullReq); err != nil {
		return nil, err
	}
	log.Printf("PullConvList request: %+v", pullReq)
	resp, err := s.messageService.PullConvList(ctx, pullReq)
	if err != nil {
		return nil, err
	}
	log.Printf("PullConvList response: %+v", resp)
	return resp, nil
}
