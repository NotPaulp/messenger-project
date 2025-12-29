package posts

import (
	"encoding/json"
	common "messenger-project/internal/handlers/common"
	"messenger-project/internal/kafka"
	"messenger-project/internal/models"
	"net/http"
	"time"
)

type PostHandler struct {
	producer *kafka.Producer
}

func NewPostHandler(producer *kafka.Producer) *PostHandler {
	return &PostHandler{
		producer: producer,
	}
}

func (h *PostHandler) Publish(w http.ResponseWriter, r *http.Request) {
	authorUsername, ok := r.Context().Value("username").(string)
	if !ok || authorUsername == "" {
		http.Error(w, "Unauthorized: username not found in token", http.StatusUnauthorized)
		return
	}
	if !common.CheckHTTPMethod(w, r, http.MethodPost) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	var req models.PublishPostRequest
	if !common.DecodeRequest(w, r, &req) {
		return
	}

	post := models.Post{
		ID:             time.Now().Unix(),
		AuthorUsername: authorUsername,
		Body:           req.Body,
		Comments:       []models.Comment{},
		PublishedAt:    time.Now(),
	}

	if err := h.producer.ProduceMessage(r.Context(), &post); err != nil {
		http.Error(w, "Error publishing the post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Post published successfully",
		"post": models.Post{
			ID:             post.ID,
			AuthorUsername: post.AuthorUsername,
			Body:           post.Body,
			Comments:       post.Comments,
			PublishedAt:    post.PublishedAt,
		},
	})
}
