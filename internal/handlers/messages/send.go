package messages

import (
	"encoding/json"
	common "messenger-project/internal/handlers/common"
	"messenger-project/internal/kafka"
	"messenger-project/internal/models"
	"net/http"
	"time"
)

type MessageHandler struct {
	producer *kafka.Producer
}

func NewMessageHandler(producer *kafka.Producer) *MessageHandler {
	return &MessageHandler{
		producer: producer,
	}
}

func (h *MessageHandler) SendHandler(w http.ResponseWriter, r *http.Request) {
	senderUsername, ok := r.Context().Value("username").(string)
	if !ok || senderUsername == "" {
		http.Error(w, "Unauthorized: username not found in token", http.StatusUnauthorized)
		return
	}
	if !common.CheckHTTPMethod(w, r, http.MethodPost) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	var req models.SendMessageRequest
	if !common.DecodeRequest(w, r, &req) {
		return
	}

	msg := models.Message{
		ID:               time.Now().Unix(),
		SenderUsername:   senderUsername,
		ReceiverUsername: req.ReceiverUsername,
		Body:             req.Body,
		SentAt:           time.Now(),
		Status:           0,
		StatusUpdatedAt:  time.Now(),
		Retries:          0,
	}

	if err := h.producer.ProduceMessage(r.Context(), &msg); err != nil {
		http.Error(w, "Error sending the message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Message sent successfully",
		"msg": models.Message{
			ID:               msg.ID,
			SenderUsername:   msg.SenderUsername,
			ReceiverUsername: msg.ReceiverUsername,
			Body:             msg.Body,
			SentAt:           msg.SentAt,
		},
	})
}
