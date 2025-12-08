package pusher

import (
	"backend/internal/im"
	"backend/internal/model"
	"backend/internal/pkg/constant"
	"backend/internal/pkg/database"
	"backend/internal/pkg/kafka"
	"backend/internal/service"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Pusher struct {
	wsServer *im.WsServer
	group    *service.GroupService
}

func InitAndRun(wsServer *im.WsServer) {
	pusher := Pusher{
		wsServer: wsServer,
		group:    service.NewGroupService(database.GetDB()),
	}
	go pusher.PushMessageToUser()
}

func (p *Pusher) PushMessageToUser() error {
	group, err := kafka.NewConsumerGroup(kafka.OnlinePushGroupID)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	defer group.Close()
	go func() {
		for err := range group.Errors() {
			log.Printf("ERROR: %v", err)
		}
	}()
	pushToUsers := func(msg *model.Message) error {
		log.Printf("Push message to users: %+v", msg)
		switch msg.ConvType {
		case constant.SingleChatType:
			clients, have := p.wsServer.Clients.GetAll(msg.TargetID)
			if !have {
				return nil
			}
			for _, client := range clients {
				if err := client.PushMessage(context.Background(), msg); err != nil {
					log.Printf("push message to user %s failed: %v", msg.TargetID, err)
				}
			}
		case constant.GroupChatType:
			log.Printf("[push] group message push not implemented")
			memberInfos, _ := p.group.GetGroupMemberList(context.Background(), msg.TargetID)
			for _, member := range memberInfos {
				memberID := member.UserID
				clients, have := p.wsServer.Clients.GetAll(memberID)
				if !have {
					continue
				}
				for _, client := range clients {
					if err := client.PushMessage(context.Background(), msg); err != nil {
						log.Printf("push message to user %s failed: %v", memberID, err)
					}
				}
			}
		}
		return nil
	}
	for {
		err := group.Consume(context.Background(), []string{kafka.OnlinePushTopic}, onlinePushHandler{fn: pushToUsers})
		if err != nil {
			panic(err)
		}
	}
}

type onlinePushHandler struct {
	fn func(msg *model.Message) error
}

func (onlinePushHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (onlinePushHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (p onlinePushHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var m model.Message
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Printf("push: invalid message json topic=%s partition=%d offset=%d err=%v", msg.Topic, msg.Partition, msg.Offset, err)
			// 解析失败，决定是否 mark（通常先记录并 mark 防止死循环），或存储以便人工检查
			sess.MarkMessage(msg, "")
			continue
		}

		// 在这里同步处理消息（例如推到在线用户）
		if err := p.fn(&m); err != nil {
			log.Printf("push: handle message failed topic=%s offset=%d err=%v", msg.Topic, msg.Offset, err)
			// 根据策略选择是否重试。此处仍 mark，避免阻塞其他消息
		}

		// 处理成功后标记偏移
		sess.MarkMessage(msg, "")
	}
	return nil
}
