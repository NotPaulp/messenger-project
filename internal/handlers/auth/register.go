package auth

import (
	"encoding/json"
	"messenger-project/internal/auth"
	"messenger-project/internal/models"
	"messenger-project/internal/repository"
	"messenger-project/pkg/utils"
	"net/http"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTTPMethod(w, r, http.MethodPost) {
		return
	}
	if !checkHTTPContentType(w, r, "application/json") {
		return
	}

	var req models.RegisterRequest
	if !decodeRequest(w, r, &req) {
		return
	}
	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	vulnerabilityMessage := utils.CheckPasswordVulnerabilities(req.Password)
	if vulnerabilityMessage != "" {
		http.Error(w, vulnerabilityMessage, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:        "user-" + time.Now().Format("20060102150405"),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repository.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User registered successfully",
		"user": models.RegisterResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
