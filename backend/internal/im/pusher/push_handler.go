package pusher

import (
	"backend/internal/model"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

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
