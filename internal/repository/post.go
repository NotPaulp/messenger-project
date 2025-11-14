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

func CreatePost(post *models.Post) error {
	_, err := GetUserByUsername(post.AuthorUsername)
	if err != nil {
		return err
	}

	if database.PostsCollection == nil {
		return fmt.Errorf("posts collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.PostsCollection.InsertOne(ctx, post)
	if err != nil {
		return fmt.Errorf("insert post: %w", err)
	}
	return nil
}

func GetLastPost(authorUsername string) (*models.Post, error) {
	if database.PostsCollection == nil {
		return nil, fmt.Errorf("posts collection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"author_username": authorUsername,
	}

	opts := options.FindOne().SetSort(bson.D{{"sent_at", -1}})
	var post models.Post
	err := database.PostsCollection.FindOne(ctx, filter, opts).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no posts from %s", authorUsername)
		}
		return nil, err
	}
	return &post, nil
}

func GetAllPosts(authorUsername string) ([]models.Post, error) {
	if database.PostsCollection == nil {
		return nil, fmt.Errorf("posts collection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"author_username": authorUsername,
	}
	cursor, err := database.PostsCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find posts: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.Post
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("cursor all: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].PublishedAt.Before(results[j].PublishedAt)
	})
	return results, nil
}
