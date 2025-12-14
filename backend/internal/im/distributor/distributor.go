package distributor

import (
	"backend/internal/im"
	"backend/internal/im/imrepo"
	"backend/internal/model"
	"backend/internal/pkg/batchprocessor"
	"backend/internal/pkg/cache/redis"
	"backend/internal/pkg/database"
	"backend/internal/pkg/kafka"
	"backend/internal/pkg/snowflake"
	"backend/internal/service"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Distributor struct {
	wsServer *im.WsServer
	repo     *imrepo.ImRepo
}

func NewDistributor(wsServer *im.WsServer) *Distributor {
	return &Distributor{
		wsServer: wsServer,
		repo:     imrepo.NewImRepo(database.GetDB(), redis.GetRDB()),
	}
}

func (d *Distributor) Start() {

	comingMessageConsumerGroup, err := kafka.NewConsumerGroup(kafka.ComingMessageGroupID)
	if err != nil {
		log.Println(err.Error())
	}
	defer comingMessageConsumerGroup.Close()
	batchprocessor := batchprocessor.NewBatchProcessor[*service.SendMessageReq]()
	batchprocessor.Key = func(val *service.SendMessageReq) string {
		return service.GetConversationID(val.ConvType, val.SenderID, val.TargetID)
	}

	onlinePushProducer, _ := kafka.NewSyncProducer()

	batchprocessor.Do = func(ctx context.Context, channelID int, msgs []*service.SendMessageReq) {
		// TODO:
		/*
			1. 存储消息到缓存
			2. 发送消息给在线用户
			3. 存储消息到数据库
		*/
		convID := service.GetConversationID(msgs[0].ConvType, msgs[0].SenderID, msgs[0].TargetID)
		var msgsToStore []*model.Message
		for _, msgReq := range msgs {
			msg := &model.Message{
				ID:             snowflake.GenID(),
				ConversationID: convID,
				SenderID:       msgReq.SenderID,
				MsgType:        msgReq.MsgType,
				Content:        msgReq.Content,
				ConvType:       msgReq.ConvType,
				TargetID:       msgReq.TargetID,
				SendTime:       time.Now().UnixMilli(),
			}
			msgsToStore = append(msgsToStore, msg)
		}
		// 1. 存储消息到缓存
		isNewConversation, err := d.repo.BatchStoreMsgToRedis(ctx, convID, msgsToStore)
		if err != nil {
			log.Printf("distributor: BatchStoreMsgToRedis error: %v", err)
		}

		if isNewConversation {
			log.Printf("distributor: new conversation created: %s", convID)
			d.repo.CreateConversations(ctx, service.InitConversationReq{
				SenderID: msgs[0].SenderID,
				ConvType: msgs[0].ConvType,
				TargetID: msgs[0].TargetID,
			})
		}

		// 2. 存储消息到数据库
		go d.repo.BatchStoreMsgToDB(context.Background(), msgsToStore)

		// 2. 发送消息给在线用户
		for _, msg := range msgsToStore {
			data, _ := json.Marshal(msg)
			onlinePushProducer.SendMessage(&sarama.ProducerMessage{
				Topic: kafka.OnlinePushTopic,
				Value: sarama.ByteEncoder(data),
			})
		}
	}

	consumeMsgHandler := func(msg *service.SendMessageReq) error {
		log.Printf("distributor: received message to distribute: %+v", msg)
		batchprocessor.Enqueue(msg)
		return nil
	}

	go func() {
		for {
			// 必须在循环中调用 Consume
			err := comingMessageConsumerGroup.Consume(context.Background(), []string{kafka.ComingMessageTopic}, &msgHandler{fn: consumeMsgHandler})
			if err != nil {
				log.Printf("distributor: consumer group error: %v", err)
				// 避免错误导致死循环空转，稍微休眠一下
				time.Sleep(time.Second)
			}
			// 如果 Context 结束了，退出循环
			if context.Background().Err() != nil {
				return
			}
		}
	}()
	batchprocessor.Start()
}

type msgHandler struct {
	fn func(*service.SendMessageReq) error
}

func NewMsgHandler(fn func(*service.SendMessageReq) error) *msgHandler {
	return &msgHandler{fn: fn}
}

func (h *msgHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *msgHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *msgHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var m service.SendMessageReq
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Printf("distributor: invalid message json topic=%s partition=%d offset=%d err=%v", msg.Topic, msg.Partition, msg.Offset, err)
			// mark to avoid reprocessing invalid payloads
			sess.MarkMessage(msg, "")
			continue
		}

		if h.fn != nil {
			if err := h.fn(&m); err != nil {
				log.Printf("distributor: handler fn error: %v", err)
			}
		}

		// mark as processed
		sess.MarkMessage(msg, "")
	}
	return nil
}
