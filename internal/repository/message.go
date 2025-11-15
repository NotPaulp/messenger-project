package repository

import (
	"context"
	"fmt"
	"messenger-project/internal/database"
	"messenger-project/internal/models"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMessage(msg *models.Message) error {
	_, err := GetUserByUsername(msg.ReceiverUsername)
	if err != nil {
		return err
	}

	if database.MessagesCollection == nil {
		return fmt.Errorf("messages collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.MessagesCollection.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

// PostgreSQL
// func CreateMessage(msg *models.Message) error {
// 	msgMap := map[string]any{
// 		"id":                msg.ID,
// 		"sender_username":   msg.SenderUsername,
// 		"receiver_username": msg.ReceiverUsername,
// 		"body":              msg.Body,
// 		"sent_at":           msg.SentAt,
// 	}
// 	_, err := GetUserByUsername(msg.ReceiverUsername)
// 	if err != nil {
// 		return err

// 	}
// 	err = database.Create(database.DB, "messages", msgMap)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func GetLastMessage(senderUsername, receiverUsername string) (*models.Message, error) {
	if database.MessagesCollection == nil {
		return nil, fmt.Errorf("messages collection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"sender_username":   senderUsername,
		"receiver_username": receiverUsername,
	}

	opts := options.FindOne().SetSort(bson.D{{"sent_at", -1}})
	var msg models.Message
	err := database.MessagesCollection.FindOne(ctx, filter, opts).Decode(&msg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no messages from %s", senderUsername)
		}
		return nil, err
	}
	return &msg, nil
}

// PostgreSQL
// func GetLastMessage(senderUsername, receiverUsername string) (*models.Message, error) {
// 	msgsData, err := database.Read(database.DB, "messages", []string{"sender_username", "receiver_username"}, []any{senderUsername, receiverUsername})

// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(msgsData) == 0 {
// 		return nil, fmt.Errorf("No messages from: %s", senderUsername)
// 	}

// 	sort.Slice(msgsData, func(i, j int) bool {
// 		return msgsData[i]["sent_at"].(time.Time).Before(msgsData[j]["sent_at"].(time.Time))
// 	})

// 	lastMsgData := msgsData[len(msgsData)-1]
// 	msg := &models.Message{
// 		ID:               lastMsgData["id"].(int64),
// 		SenderUsername:   lastMsgData["sender_username"].(string),
// 		ReceiverUsername: lastMsgData["receiver_username"].(string),
// 		Body:             lastMsgData["body"].(string),
// 		SentAt:           lastMsgData["sent_at"].(time.Time),
// 	}
// 	return msg, nil
// }

func GetAllMessages(senderUsername, receiverUsername string) ([]models.Message, error) {
	if database.MessagesCollection == nil {
		return nil, fmt.Errorf("messages collection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"sender_username":   senderUsername,
		"receiver_username": receiverUsername,
	}
	cursor, err := database.MessagesCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.Message
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("cursor all: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].SentAt.Before(results[j].SentAt)
	})
	return results, nil
}

// PostgreSQL
// func GetAllMessages(senderUsername, receiverUsername string) ([]models.Message, error) {
// 	msgsData, err := database.Read(database.DB, "messages", []string{"sender_username", "receiver_username"}, []any{senderUsername, receiverUsername})

// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(msgsData) == 0 {
// 		return nil, fmt.Errorf("No messages from: %s", senderUsername)
// 	}
// 	sort.Slice(msgsData, func(i, j int) bool {
// 		return msgsData[i]["sent_at"].(time.Time).Before(msgsData[j]["sent_at"].(time.Time))
// 	})
// 	var msgs []models.Message
// 	for _, msgData := range msgsData {
// 		msg := models.Message{
// 			ID:               msgData["id"].(int64),
// 			SenderUsername:   msgData["sender_username"].(string),
// 			ReceiverUsername: msgData["receiver_username"].(string),
// 			Body:             msgData["body"].(string),
// 			SentAt:           msgData["sent_at"].(time.Time),
// 		}
// 		msgs = append(msgs, msg)
// 	}
// 	return msgs, nil
// }

// Add to internal/repository/message.go
func DeleteMessage(messageID int64, username string) error {
	if database.MessagesCollection == nil {
		return fmt.Errorf("messages collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"id":              messageID,
		"sender_username": username,
	}

	result, err := database.MessagesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("message not found or you don't have permission to delete it")
	}

	return nil
}
