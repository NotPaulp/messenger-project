package models

import (
	"time"
)

type Post struct {
	ID             int64     `json:"id" bson:"id"`
	AuthorUsername string    `json:"author_username" bson:"author_username"`
	Body           string    `json:"body" bson:"body"`
	Comments       []Comment `json:"comments" bson:"comments"`
	PublishedAt    time.Time `json:"published_at" bson:"published_at"`
	Categories     []string  `json:"categories" bson:"categories"`
	AnalyzedAt     time.Time `json:"analyzed_at" bson:"analyzed_at"`
}

func (p *Post) GetKafkaKey() string {
	return p.AuthorUsername
}

type PublishPostRequest struct {
	Body string `json:"body"`
}

type GetPostRequest struct {
	AuthorUsername string `json:"author_username"`
	All            bool   `json:"all"`
}

type DeletePostRequest struct {
	PostID int64 `json:"post_id"`
}
