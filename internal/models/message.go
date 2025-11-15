package models

import (
	"time"
)

type Message struct {
	ID               int64     `json:"id" bson:"id"`
	SenderUsername   string    `json:"sender_username" bson:"sender_username"`
	ReceiverUsername string    `json:"receiver_username" bson:"receiver_username"`
	Body             string    `json:"body" bson:"body"`
	SentAt           time.Time `json:"sent_at" bson:"sent_at"`
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
