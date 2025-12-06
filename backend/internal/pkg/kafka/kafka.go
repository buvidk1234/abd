package kafka

import "github.com/IBM/sarama"

var conf Config

func Init(cfg Config) {
	conf = cfg
}

func NewSyncProducer() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducer(conf.Addr, nil)
}

func NewConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(conf.Addr, groupID, nil)
}
