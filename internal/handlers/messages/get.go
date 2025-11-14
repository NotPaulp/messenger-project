package messages

import (
	"encoding/json"
	common "messenger-project/internal/handlers/common"
	"messenger-project/internal/models"
	"messenger-project/internal/repository"
	"net/http"
	"strings"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	receiverUsername, ok := r.Context().Value("username").(string)
	if !ok || receiverUsername == "" {
		http.Error(w, "Unauthorized: username not found in token", http.StatusUnauthorized)
		return
	}

	if !common.CheckHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if !common.CheckHTTPContentType(w, r, "application/json") {
		return
	}

	senderUsername := r.URL.Query().Get("sender")
	allStr := r.URL.Query().Get("all")
	if senderUsername == "" {
		http.Error(w, "sender query param is required", http.StatusBadRequest)
		return
	}

	all := strings.ToLower(allStr) == "true"
	if all {
		msgs, err := repository.GetAllMessages(senderUsername, receiverUsername)
		if err != nil {
			http.Error(w, "Error retrieving messages: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Messages retrieved successfully",
			"msg":     msgs,
		})
	}
	msg, err := repository.GetLastMessage(senderUsername, receiverUsername)
	if err != nil {
		http.Error(w, "Error retrieving the last message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Last message retrieved successfully",
		"msg": []models.Message{{
			ID:               msg.ID,
			SenderUsername:   msg.SenderUsername,
			ReceiverUsername: msg.ReceiverUsername,
			Body:             msg.Body,
			SentAt:           msg.SentAt,
		}},
	})
}
