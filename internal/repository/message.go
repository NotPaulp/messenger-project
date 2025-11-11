package repository

import (
	"fmt"
	"messenger-project/internal/database"
	"messenger-project/internal/models"
	"sort"
	"time"
)

func CreateMessage(msg *models.Message) error {
	msgMap := map[string]any{
		"id":                msg.ID,
		"sender_username":   msg.SenderUsername,
		"receiver_username": msg.ReceiverUsername,
		"body":              msg.Body,
		"sent_at":           msg.SentAt,
	}
	_, err := GetUserByUsername(msg.ReceiverUsername)
	if err != nil {
		return err

	}
	err = database.Create(database.DB, "messages", msgMap)
	if err != nil {
		return err
	}
	return nil
}

func GetLastMessage(senderUsername, receiverUsername string) (*models.Message, error) {
	msgsData, err := database.Read(database.DB, "messages", []string{"sender_username", "receiver_username"}, []any{senderUsername, receiverUsername})

	if err != nil {
		return nil, err
	}
	if len(msgsData) == 0 {
		return nil, fmt.Errorf("No messages from: %s", senderUsername)
	}

	sort.Slice(msgsData, func(i, j int) bool {
		return msgsData[i]["sent_at"].(time.Time).Before(msgsData[j]["sent_at"].(time.Time))
	})

	lastMsgData := msgsData[len(msgsData)-1]
	msg := &models.Message{
		ID:               lastMsgData["id"].(int64),
		SenderUsername:   lastMsgData["sender_username"].(string),
		ReceiverUsername: lastMsgData["receiver_username"].(string),
		Body:             lastMsgData["body"].(string),
		SentAt:           lastMsgData["sent_at"].(time.Time),
	}
	return msg, nil
}

func GetAllMessages(senderUsername, receiverUsername string) ([]models.Message, error) {
	msgsData, err := database.Read(database.DB, "messages", []string{"sender_username", "receiver_username"}, []any{senderUsername, receiverUsername})

	if err != nil {
		return nil, err
	}
	if len(msgsData) == 0 {
		return nil, fmt.Errorf("No messages from: %s", senderUsername)
	}
	sort.Slice(msgsData, func(i, j int) bool {
		return msgsData[i]["sent_at"].(time.Time).Before(msgsData[j]["sent_at"].(time.Time))
	})
	var msgs []models.Message
	for _, msgData := range msgsData {
		msg := models.Message{
			ID:               msgData["id"].(int64),
			SenderUsername:   msgData["sender_username"].(string),
			ReceiverUsername: msgData["receiver_username"].(string),
			Body:             msgData["body"].(string),
			SentAt:           msgData["sent_at"].(time.Time),
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
