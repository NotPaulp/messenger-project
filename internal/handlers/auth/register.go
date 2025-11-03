package auth

import (
	"encoding/json"
	"messenger-project/internal/auth"
	"messenger-project/internal/models"
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
	if !validateRegisterRequest(w, &req) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user": models.RegisterResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

func checkHTTPMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
func checkHTTPContentType(w http.ResponseWriter, r *http.Request, contentType string) bool {
	if r.Header.Get("Content-Type") != contentType {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

func decodeRequest(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
func validateRegisterRequest(w http.ResponseWriter, req *models.RegisterRequest) bool {

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return false
	}
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return false
	}
	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return false
	}
	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return false
	}
	return true
}
