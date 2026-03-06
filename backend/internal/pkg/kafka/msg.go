package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type KafkaMessage struct {
	Ctx     context.Context
	Msg     *sarama.ConsumerMessage
	Session sarama.ConsumerGroupSession
}
