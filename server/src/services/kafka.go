package services

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "maintenance-task-events",
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}, nil
}

func (kp *KafkaProducer) SendMessage(message []byte) error {
	err := kp.writer.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})
	if err != nil {
		log.Println("Failed to send message to Kafka:", err)
		return err
	}

	log.Println("Message sent to Kafka successfully")
	return nil
}

func (kp *KafkaProducer) Close() error {
	err := kp.writer.Close()
	if err != nil {
		log.Println("Failed to close Kafka producer:", err)
		return err
	}
	log.Println("Kafka producer closed successfully")
	return nil
}
