package auth

import (
	"encoding/json"
	"messenger-project/internal/models"
	"messenger-project/internal/redis"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTTPMethod(w, r, http.MethodPost) {
		return
	}
	if !checkHTTPContentType(w, r, "application/json") {
		return
	}

	var req models.LogoutRequest
	if !decodeRequest(w, r, &req) {
		return
	}

	if req.Token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	err := redis.BlacklistJWT(req.Token)
	if err != nil {
		http.Error(w, "Error blacklisting JWT token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "Logged out"})
}
