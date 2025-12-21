package comments

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	common "messenger-project/internal/handlers/common"
	comments "messenger-project/internal/repository/api-gateway"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("username").(string)
	if !ok {
		http.Error(w, "Unauthorized: username not found in token", http.StatusUnauthorized)
		return
	}

	if !common.CheckHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "post_id query param is required", http.StatusBadRequest)
		return
	}
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "post_id must be an integer", http.StatusBadRequest)
		return
	}

	allStr := r.URL.Query().Get("all")
	all := strings.ToLower(allStr) == "true"

	if all {
		comments, err := comments.GetComments(postID)
		if err != nil {
			http.Error(w, "Error retrieving comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message":  "Comments retrieved successfully",
			"comments": comments,
		})
		return
	}

	comment, err := comments.GetLastComment(postID)
	if err != nil {
		http.Error(w, "Error retrieving last comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Last comment retrieved successfully",
		"comment": comment,
	})
}
