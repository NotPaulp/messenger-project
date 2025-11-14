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
}

type PublishPostRequest struct {
	Body string `json:"body"`
}

type GetPostRequest struct {
	AuthorUsername string `json:"author_username"`
	All            bool   `json:"all"`
}
