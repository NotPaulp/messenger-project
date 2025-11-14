package database

import (
	"context"
	"fmt"
	"messenger-project/pkg/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MessagesCollection *mongo.Collection
var PostsCollection *mongo.Collection

func InitMongo(ctx context.Context) error {
	cfg := config.Load()
	mongoURL := cfg.MONGO_URL
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}

	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()
	if err := client.Ping(ctxPing, nil); err != nil {
		_ = client.Disconnect(ctx)
		return fmt.Errorf("mongo ping: %w", err)
	}

	MongoClient = client

	dbName := cfg.MONGO_DB
	if dbName == "" {
		dbName = "messenger"
	}

	messagesCollection := cfg.MONGO_MESSAGES_COLLECTION
	if messagesCollection == "" {
		messagesCollection = "messages"
	}

	MessagesCollection = client.Database(dbName).Collection(messagesCollection)

	postsCollection := cfg.MONGO_POSTS_COLLECTION
	if postsCollection == "" {
		postsCollection = "posts"
	}

	PostsCollection = client.Database(dbName).Collection(postsCollection)

	return nil
}

func DisconnectMongo(ctx context.Context) error {
	if MongoClient == nil {
		return nil
	}
	return MongoClient.Disconnect(ctx)
}
