package comments

import (
	"encoding/json"
	"net/http"
	"time"

	common "messenger-project/internal/handlers/common"
	"messenger-project/internal/models"
	comments "messenger-project/internal/repository/api-gateway"
)

func PublishCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	var req models.PublishCommentRequest
	if !common.DecodeRequest(w, r, &req) {
		return
	}

	comment := &models.Comment{
		ID:             time.Now().Unix(),
		AuthorUsername: authorUsername,
		Body:           req.Body,
		PublishedAt:    time.Now(),
	}

	var err error
	if comment, err = comments.AddCommentAndReturn(req.PostID, comment); err != nil {
		http.Error(w, "Error adding comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Comment published successfully",
		"comment": comment,
	})
}
