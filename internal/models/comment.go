package models

import (
	"time"
)

type Comment struct {
	ID             int64     `json:"id" bson:"id"`
	AuthorUsername string    `json:"author_username" bson:"author_username"`
	Body           string    `json:"body" bson:"body"`
	PublishedAt    time.Time `json:"published_at" bson:"published_at"`
}

type PublishCommentRequest struct {
	PostID int64  `json:"post_id"`
	Body   string `json:"body"`
}

type GetCommentRequest struct {
	PostID         int64  `json:"post_id"`
	AuthorUsername string `json:"author_username"`
	All            bool   `json:"all"`
}

type DeleteCommentRequest struct {
	PostID    int64 `json:"post_id"`
	CommentID int64 `json:"comment_id"`
}
