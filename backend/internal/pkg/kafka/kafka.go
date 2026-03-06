package kafka

import (
	"bytes"
	"strings"

	"github.com/IBM/sarama"
)

var conf Config

func Init(cfg Config) {
	conf = cfg
}

func NewSyncProducer() (sarama.SyncProducer, error) {
	kfk := sarama.NewConfig()
	kfk.Producer.Return.Successes = true
	kfk.Producer.Return.Errors = true
	kfk.Producer.Partitioner = sarama.NewHashPartitioner
	switch strings.ToLower(conf.ProducerAck) {
	case "no_response":
		kfk.Producer.RequiredAcks = sarama.NoResponse
	case "wait_for_local":
		kfk.Producer.RequiredAcks = sarama.WaitForLocal
	case "wait_for_all":
		kfk.Producer.RequiredAcks = sarama.WaitForAll
	default:
		kfk.Producer.RequiredAcks = sarama.WaitForAll
	}
	if conf.CompressType == "" {
		kfk.Producer.Compression = sarama.CompressionNone
	} else {
		if err := kfk.Producer.Compression.UnmarshalText(bytes.ToLower([]byte(conf.CompressType))); err != nil {
			return nil, err
		}
	}
	return sarama.NewSyncProducer(conf.Addr, kfk)
}

func NewConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	kfk := sarama.NewConfig()
	kfk.Consumer.Offsets.Initial = sarama.OffsetNewest
	// kfk.Consumer.Offsets.AutoCommit.Enable = false
	return sarama.NewConsumerGroup(conf.Addr, groupID, kfk)
}
