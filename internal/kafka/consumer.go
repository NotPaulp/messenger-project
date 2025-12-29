package kafka

import (
	"context"
	"encoding/json"
	"log"
	"messenger-project/internal/models"
	api_gateway "messenger-project/internal/repository/api-gateway"
	"messenger-project/internal/websocket"
	"messenger-project/pkg/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	websocketHub *websocket.Hub
}

func NewConsumer(cfg *config.Config, hub *websocket.Hub) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KAFKA_BROKERS},
		Topic:   cfg.KAFKA_TOPIC_CONTENT,
		GroupID: "message-consumer-group",
	})

	return &Consumer{
		reader:       reader,
		websocketHub: hub,
	}
}

func (c *Consumer) Start(ctx context.Context, mlanalyzerproducer *Producer) {
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
			contentType := ""
			for _, h := range msg.Headers {
				if h.Key == "content_type" {
					contentType = string(h.Value)
					break
				}
			}
			switch contentType {
			case "message":
				var message models.Message
				if err := json.Unmarshal(msg.Value, &message); err != nil {
					log.Printf("Error unmarshaling kafka message: %v", err)
					continue
				}

				if err := api_gateway.CreateMessage(&message); err != nil {
					log.Printf("Error saving message to database: %v", err)
					continue
				}

				err := mlanalyzerproducer.ProduceMessage(ctx, &message)
				if err != nil {
					log.Printf("Error producing message to ML analyzer: %v", err)
				}
				log.Printf("Message processed and saved: %s -> %s",
					message.SenderUsername, message.ReceiverUsername)

				if c.websocketHub == nil {
					continue
				}

				if err := c.websocketHub.BroadcastToUser(message.ReceiverUsername, message); err != nil {
					log.Printf("Error broadcasting message to user %s: %v",
						message.ReceiverUsername, err)
					continue
				}

				toUpdate := map[string]any{
					"status":            1,
					"status_updated_at": time.Now(),
				}

				if err := api_gateway.UpdateMessageByID(message.ID, toUpdate); err != nil {
					log.Printf("Error updating message %d status to sent in database: %v",
						message.ID, err)
				}
			case "post":
				var post models.Post

				if err := json.Unmarshal(msg.Value, &post); err != nil {
					log.Printf("Error unmarshaling kafka message: %v", err)
					continue
				}

				if err := api_gateway.CreatePost(&post); err != nil {
					log.Printf("Error saving post to database: %v", err)
					continue
				}

				err := mlanalyzerproducer.ProduceMessage(ctx, &post)
				if err != nil {
					log.Printf("Error producing post to ML analyzer: %v", err)
				}
				log.Printf("Post processed and saved")
			case "comment":
				var comment models.Comment
				if err := json.Unmarshal(msg.Value, &comment); err != nil {
					log.Printf("Error unmarshaling kafka message: %v", err)
					continue
				}
				if _, err := api_gateway.AddCommentAndReturn(&comment); err != nil {
					log.Printf("Error saving comment to database: %v", err)
					continue
				}

				err := mlanalyzerproducer.ProduceMessage(ctx, &comment)
				if err != nil {
					log.Printf("Error producing comment to ML analyzer: %v", err)
				}
				log.Printf("Comment processed and saved")

			}

		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
