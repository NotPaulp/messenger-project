package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"messenger-project/internal/models"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}

	return &Producer{
		writer: writer,
	}
}

func (producer *Producer) ProduceMessage(ctx context.Context, message *models.Message) error {
	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(message.ReceiverUsername),
		Value: value,
		Headers: []kafka.Header{
			{
				Key:   "sender_username",
				Value: []byte(message.SenderUsername),
			},
		},
	}

	err = producer.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
