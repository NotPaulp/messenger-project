package comments

import (
	"encoding/json"
	common "messenger-project/internal/handlers/common"
	"messenger-project/internal/models"
	comments "messenger-project/internal/repository/api-gateway"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized: username not found in token", http.StatusUnauthorized)
		return
	}

	if !common.CheckHTTPMethod(w, r, http.MethodDelete) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	var req models.DeleteCommentRequest
	if !common.DecodeRequest(w, r, &req) {
		return
	}

	if req.PostID == 0 {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	if req.CommentID == 0 {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	err := comments.DeleteComment(req.PostID, req.CommentID, username)
	if err != nil {
		http.Error(w, "Error deleting comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":    "Comment deleted successfully",
		"post_id":    req.PostID,
		"comment_id": req.CommentID,
	})
}
