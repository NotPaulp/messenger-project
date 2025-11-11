package models

import (
	"time"
)

type Message struct {
	ID               int64     `json:"id" db:"id"`
	SenderUsername   string    `json:"sender_username" db:"sender_username"`
	ReceiverUsername string    `json:"receiver_username" db:"receiver_username"`
	Body             string    `json:"body" db:"body"`
	SentAt           time.Time `json:"sent_at" db:"sent_at"`
}

type SendMessageRequest struct {
	ReceiverUsername string `json:"receiver_username"`
	Body             string `json:"body"`
}

type GetMessageRequest struct {
	SenderUsername string `json:"sender_username"`
	All            bool   `json:"all"`
}
