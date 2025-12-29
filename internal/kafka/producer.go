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

type Content interface {
	GetKafkaKey() string
}

func (producer *Producer) ProduceMessage(ctx context.Context, content Content) error {
	value, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	contentType := getContentType(content)
	msg := kafka.Message{
		Key:   []byte(content.GetKafkaKey()),
		Value: value,
		Headers: []kafka.Header{
			{
				Key:   "content_type",
				Value: []byte(contentType),
			},
		},
	}

	if err := producer.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func getContentType(c Content) string {
	switch any(c).(type) {
	case *models.Message:
		return "message"
	case *models.Post:
		return "post"
	case *models.Comment:
		return "comment"
	default:
		return "unknown"
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
