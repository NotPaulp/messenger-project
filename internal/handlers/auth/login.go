package auth

import (
	"encoding/json"
	"messenger-project/internal/auth"
	"messenger-project/internal/handlers/common"
	"messenger-project/internal/models"
	users "messenger-project/internal/repository/user-service"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !common.CheckHTTPMethod(w, r, http.MethodPost) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	var req models.LoginRequest
	if !common.DecodeRequest(w, r, &req) {
		return
	}
	if !validateLoginRequest(w, &req) {
		return
	}

	user, err := users.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Error while getting the user out of the database", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User authenticated successfully",
		"token":   token,
		"user": models.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

func validateLoginRequest(w http.ResponseWriter, req *models.LoginRequest) bool {

	if req.Username == "" && req.Email == "" {
		http.Error(w, "Username or email is required", http.StatusBadRequest)
		return false
	}
	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return false
	}
	return true
}
