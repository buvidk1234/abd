package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating cluster admin: %v", err)
	}
	defer admin.Close()

	topics := []string{"coming_message_topic", "online_push_topic"}

	for _, topic := range topics {
		log.Printf("Deleting topic: %s", topic)
		err := admin.DeleteTopic(topic)
		if err != nil {
			log.Printf("Error deleting topic %s: %v (it might not exist)", topic, err)
		}
	}

	// Wait for deletion to propagate
	time.Sleep(2 * time.Second)

	for _, topic := range topics {
		log.Printf("Creating topic: %s", topic)
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			log.Printf("Error creating topic %s: %v", topic, err)
		}
	}

	log.Println("Kafka topics reset successfully.")
}
