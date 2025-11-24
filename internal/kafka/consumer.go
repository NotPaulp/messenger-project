package kafka

import (
	"context"
	"encoding/json"
	"log"
	"messenger-project/internal/models"
	"messenger-project/internal/repository"
	"messenger-project/internal/websocket"
	"messenger-project/pkg/config"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	websocketHub *websocket.Hub
}

func NewConsumer(cfg *config.Config, hub *websocket.Hub) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KAFKA_BROKERS},
		Topic:   cfg.KAFKA_TOPIC_MESSAGES,
		GroupID: "message-consumer-group",
	})

	return &Consumer{
		reader:       reader,
		websocketHub: hub,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("Error reading kafka message: %v", err)
				continue
			}

			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Error unmarshaling kafka message: %v", err)
				continue
			}

			if err := repository.CreateMessage(&message); err != nil {
				log.Printf("Error saving message to database: %v", err)
				continue
			}

			log.Printf("Message processed and saved: %s -> %s", message.SenderUsername, message.ReceiverUsername)
			if c.websocketHub != nil {
				c.websocketHub.BroadcastToUser(message.ReceiverUsername, message)
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
