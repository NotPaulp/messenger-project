package messages

import (
	"encoding/json"
	common "messenger-project/internal/handlers/common"
	"messenger-project/internal/models"
	messages "messenger-project/internal/repository/api-gateway"
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

	var req models.DeleteMessageRequest
	if !common.DecodeRequest(w, r, &req) {
		return
	}

	if req.MessageID == 0 {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	err := messages.DeleteMessage(req.MessageID, username)
	if err != nil {
		http.Error(w, "Error deleting message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":    "Message deleted successfully",
		"message_id": req.MessageID,
	})
}
