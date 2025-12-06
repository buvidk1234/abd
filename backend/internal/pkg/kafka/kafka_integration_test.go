package kafka

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

// This is a simple integration test that connects to a Kafka broker at 192.168.6.130:9092,
// sends one message and verifies a consumer group receives it. Adjust the broker address
// or topic as needed. The test expects the broker to allow topic auto-creation or the
// topic to already exist.
func TestSendAndReceive(t *testing.T) {
	// broker address
	broker := "192.168.6.130:9092"

	// quick reachability check so test skips fast when broker is not available
	connOK := func(addr string, timeout time.Duration) bool {
		d := net.Dialer{Timeout: timeout}
		c, err := d.DialContext(context.Background(), "tcp", addr)
		if err != nil {
			return false
		}
		_ = c.Close()
		return true
	}
	if !connOK(broker, 1*time.Second) {
		t.Skipf("kafka broker %s unreachable, skipping integration test", broker)
	}

	// configure package
	Init(Config{
		Addr: []string{broker},
	})

	// create producer
	prod, err := NewSyncProducer()
	if err != nil {
		t.Fatalf("NewSyncProducer failed: %v", err)
	}
	defer func() { _ = prod.Close() }()

	topic := "test_topic_demo"
	groupID := fmt.Sprintf("test-group-%d", time.Now().UnixNano())

	// create consumer group with config (start from oldest to avoid missing messages)
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	group, err := sarama.NewConsumerGroup([]string{broker}, groupID, cfg)
	if err != nil {
		t.Fatalf("NewConsumerGroup failed: %v", err)
	}
	defer func() { _ = group.Close() }()

	// channel to receive one message
	msgCh := make(chan *sarama.ConsumerMessage, 1)

	handler := &simpleHandler{msgCh: msgCh}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start consumer loop
	go func() {
		for {
			if err := group.Consume(ctx, []string{topic}, handler); err != nil {
				// log and retry a bit
				t.Logf("consumer error: %v", err)
				time.Sleep(500 * time.Millisecond)
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	// Give consumer some time to join the group
	time.Sleep(1 * time.Second)

	// send a message
	body := []byte("hello kafka test")
	pm := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(body)}
	partition, offset, err := prod.SendMessage(pm)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	t.Logf("sent message partition=%d offset=%d", partition, offset)

	// wait for consumer to receive (give a bit more time)
	select {
	case m := <-msgCh:
		if string(m.Value) != string(body) {
			t.Fatalf("received payload mismatch: got=%s want=%s", string(m.Value), string(body))
		}
		t.Logf("received message from topic=%s partition=%d offset=%d", m.Topic, m.Partition, m.Offset)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for message")
	}

	cancel()
}

type simpleHandler struct {
	msgCh chan *sarama.ConsumerMessage
}

func (h *simpleHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *simpleHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *simpleHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for m := range claim.Messages() {
		// deliver and mark
		h.msgCh <- m
		sess.MarkMessage(m, "")
		return nil // stop after first message
	}
	return nil
}
