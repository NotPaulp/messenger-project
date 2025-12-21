package retry

import (
	"context"
	"log"
	"messenger-project/internal/models"
	repository "messenger-project/internal/repository/api-gateway"
	"messenger-project/internal/websocket"
	"time"
)

type RetryService struct {
	hub *websocket.Hub
}

func NewRetryService(hub *websocket.Hub) *RetryService {
	return &RetryService{hub: hub}
}

func (s *RetryService) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("RetryService: context canceled, stopping")
			return

		case <-ticker.C:
			s.processPending()
		}
	}
}

func (s *RetryService) processPending() {

	where := map[string]any{
		"status": models.StatusPending,
	}

	messages, err := repository.GetAllMessagesWhere(where)
	if err != nil {
		log.Printf("RetryService: error getting pending messages: %v", err)
		return
	}

	if len(messages) == 0 {
		return
	}

	log.Printf("RetryService: found %d pending messages to process", len(messages))

	for _, message := range messages {
		if message.Retries >= models.MaxRetries {
			toUpdate := map[string]any{
				"status":            models.StatusFailed,
				"status_updated_at": time.Now(),
			}
			if err := repository.UpdateMessageByID(message.ID, toUpdate); err != nil {
				log.Printf("RetryService: error marking message %d as failed: %v", message.ID, err)
			}
			continue
		}

		if err := s.hub.BroadcastToUser(message.ReceiverUsername, message); err != nil {
			log.Printf("RetryService: error broadcasting message %d to %s: %v",
				message.ID, message.ReceiverUsername, err)

			toUpdate := map[string]any{
				"retries": message.Retries + 1,
			}

			if message.Retries+1 >= models.MaxRetries {
				toUpdate["status"] = models.StatusFailed
				toUpdate["status_updated_at"] = time.Now()
			}

			if err := repository.UpdateMessageByID(message.ID, toUpdate); err != nil {
				log.Printf("RetryService: error updating retries for message %d: %v", message.ID, err)
			}

			continue
		}

		toUpdate := map[string]any{
			"status":            models.StatusSent,
			"sent_at":           time.Now(),
			"status_updated_at": time.Now(),
		}

		if err := repository.UpdateMessageByID(message.ID, toUpdate); err != nil {
			log.Printf("RetryService: error marking message %d as sent: %v", message.ID, err)
		}
	}
}
