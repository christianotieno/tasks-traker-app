package handlers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

// HandleKafkaMessages starts consuming Kafka messages.
func HandleKafkaMessages(brokers []string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "task-app",
		Topic:   "task-events",
	})

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Println("Failed to close Kafka reader:", err)
		}
	}()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading Kafka message:", err)
			continue
		}

		// Handle the Kafka message here
		log.Printf("Received Kafka message: %s\n", string(msg.Value))
	}
}
