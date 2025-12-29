package models

import "time"

type MessageAnalysisResult struct {
	MessageID      int64     `json:"message_id" bson:"message_id"`
	Category       string    `json:"category" bson:"category"`
	ToxicityScore  float32   `json:"toxicity_score" bson:"toxicity_score"`
	IsToxic        bool      `json:"toxic" bson:"toxic"`
	AnalyzedAt     time.Time `json:"analyzed_at" bson:"analyzed_at"`
	SenderUsername string    `json:"sender_username" bson:"sender_username"`
}

type PostAnalysisResult struct {
	PostID         int64     `json:"post_id" bson:"post_id"`
	Categories     []string  `json:"categories" bson:"categories"`
	AnalyzedAt     time.Time `json:"analyzed_at" bson:"analyzed_at"`
	AuthorUsername string    `json:"author_username" bson:"author_username"`
}

type CommentAnalysisResult struct {
	CommentID      int64     `json:"comment_id" bson:"comment_id"`
	PostID         int64     `json:"post_id" bson:"post_id"`
	Spam           bool      `json:"spam" bson:"spam"`
	ToxicityScore  float32   `json:"toxicity_score" bson:"toxicity_score"`
	IsToxic        bool      `json:"toxic" bson:"toxic"`
	AnalyzedAt     time.Time `json:"analyzed_at" bson:"analyzed_at"`
	AuthorUsername string    `json:"author_username" bson:"author_username"`
}
