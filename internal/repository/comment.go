package repository

import (
	"context"
	"fmt"
	"messenger-project/internal/database"
	"messenger-project/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddCommentAndReturn(postID int64, comment *models.Comment) (*models.Comment, error) {
	if database.PostsCollection == nil {
		return nil, fmt.Errorf("posts collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": postID}
	update := bson.M{"$push": bson.M{"comments": comment}}

	findOpts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetProjection(bson.M{"comments": bson.M{"$slice": -1}, "_id": 0})

	var result struct {
		Comments []models.Comment `bson:"comments"`
	}

	err := database.PostsCollection.FindOneAndUpdate(ctx, filter, update, findOpts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no post found with id %d", postID)
		}
		return nil, fmt.Errorf("add comment and return: %w", err)
	}

	if len(result.Comments) == 0 {
		return nil, fmt.Errorf("comment not returned after update for post %d", postID)
	}

	last := result.Comments[len(result.Comments)-1]
	return &last, nil
}

func GetComments(postID int64) ([]models.Comment, error) {
	if database.PostsCollection == nil {
		return nil, fmt.Errorf("posts collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": postID}
	projection := bson.M{"comments": 1, "_id": 0}

	var result struct {
		Comments []models.Comment `bson:"comments"`
	}

	err := database.PostsCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no post with id %d", postID)
		}
		return nil, err
	}
	return result.Comments, nil
}

func GetLastComment(postID int64) (*models.Comment, error) {
	comments, err := GetComments(postID)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, fmt.Errorf("no comments for post %d", postID)
	}

	last := comments[len(comments)-1]
	return &last, nil
}

func DeleteComment(postID, commentID int64, username string) error {
	if database.PostsCollection == nil {
		return fmt.Errorf("posts collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"id": postID,
		"comments": bson.M{
			"$elemMatch": bson.M{
				"id":              commentID,
				"author_username": username,
			},
		},
	}

	update := bson.M{
		"$pull": bson.M{
			"comments": bson.M{
				"id":              commentID,
				"author_username": username,
			},
		},
	}

	result, err := database.PostsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("comment not found or you don't have permission to delete it")
	}

	return nil
}
