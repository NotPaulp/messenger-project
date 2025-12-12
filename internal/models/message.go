package models

import (
	"time"
)

const (
	StatusFailed  = -1
	StatusPending = 0
	StatusSent    = 1
	MaxRetries    = 3
)

type Message struct {
	ID               int64     `json:"id" bson:"id"`
	SenderUsername   string    `json:"sender_username" bson:"sender_username"`
	ReceiverUsername string    `json:"receiver_username" bson:"receiver_username"`
	Body             string    `json:"body" bson:"body"`
	Category         string    `json:"category" bson:"category"`
	Toxic            bool      `json:"toxic" bson:"toxic"`
	ToxicityScore    float32   `json:"toxicity_score" bson:"toxicity_score"`
	SentAt           time.Time `json:"sent_at" bson:"sent_at"`
	Status           int       `json:"status" bson:"status"`
	StatusUpdatedAt  time.Time `json:"status_updated_at" bson:"status_updated_at"`
	Retries          int       `json:"retries" bson:"retries"`
}

type SendMessageRequest struct {
	ReceiverUsername string `json:"receiver_username"`
	Body             string `json:"body"`
}

type GetMessageRequest struct {
	SenderUsername string `json:"sender_username"`
	All            bool   `json:"all"`
}

type DeleteMessageRequest struct {
	MessageID int64 `json:"message_id"`
}
