package posts

import (
	"encoding/json"
	"net/http"
	"strings"

	common "messenger-project/internal/handlers/common"
	posts "messenger-project/internal/repository/api-gateway"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	authorUsername, ok := r.Context().Value("username").(string)
	if !ok || authorUsername == "" {
		http.Error(w, "Unauthorized: username not found in token", http.StatusUnauthorized)
		return
	}

	if !common.CheckHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	qAuthor := r.URL.Query().Get("author_username")
	if qAuthor == "" {
		qAuthor = authorUsername
	}
	allStr := r.URL.Query().Get("all")
	all := strings.ToLower(allStr) == "true"

	if all {
		posts, err := posts.GetAllPosts(qAuthor)
		if err != nil {
			http.Error(w, "Error retrieving posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Posts retrieved successfully",
			"posts":   posts,
		})
		return
	}

	post, err := posts.GetLastPost(qAuthor)
	if err != nil {
		http.Error(w, "Error retrieving last post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Last post retrieved successfully",
		"post":    post,
	})
}
