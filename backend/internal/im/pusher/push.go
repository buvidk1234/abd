package pusher

import (
	"backend/internal/im"
	"backend/internal/model"
	"backend/internal/pkg/constant"
	"backend/internal/pkg/database"
	"backend/internal/pkg/kafka"
	"backend/internal/service"
	"context"
	"log"
	"strconv"
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
			clients := make(map[int64][]*im.Client, 2)
			if targetClients, ok := p.wsServer.Clients.GetAll(msg.TargetID); ok {
				clients[msg.TargetID] = targetClients
			}
			if senderClients, ok := p.wsServer.Clients.GetAll(msg.SenderID); ok {
				// sender==target 时会自动合并到同一个 key
				clients[msg.SenderID] = append(clients[msg.SenderID], senderClients...)
			}

			// 两边都不在线就直接返回
			if len(clients) == 0 {
				return nil
			}
			ctx := context.Background()
			for uid, cs := range clients {
				for _, c := range cs {
					if err := c.PushMessage(ctx, msg); err != nil {
						log.Printf("push message to user %d failed: %v", uid, err)
					}
				}
			}
		case constant.GroupChatType:
			log.Printf("[push] group message push not implemented")
			memberInfos, _ := p.group.GetGroupMemberList(context.Background(), strconv.FormatInt(msg.TargetID, 10))
			for _, member := range memberInfos {
				memberID := member.UserID
				clients, have := p.wsServer.Clients.GetAll(memberID)
				if !have {
					continue
				}
				for _, client := range clients {
					if err := client.PushMessage(context.Background(), msg); err != nil {
						log.Printf("push message to user %d failed: %v", memberID, err)
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
