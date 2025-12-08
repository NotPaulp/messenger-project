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

func UpdateMessageStatus(messageID int64, toUpdate map[string]any) error {
	if database.MessagesCollection == nil {
		return fmt.Errorf("messages collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"id": messageID,
	}

	update := bson.M{
		"$set": toUpdate,
	}

	_, err := database.MessagesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update message: %w", err)
	}
	return nil
}

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

func GetAllMessagesWhere(where map[string]any) ([]models.Message, error) {
	if database.MessagesCollection == nil {
		return nil, fmt.Errorf("messages collection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.MessagesCollection.Find(ctx, where)
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
