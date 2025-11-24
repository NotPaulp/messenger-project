package messages

import (
	"messenger-project/internal/auth"
	"messenger-project/internal/websocket"
	"net/http"
	"strings"
)

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	tokenString := getTokenFromRequest(r)
	if tokenString == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	h.hub.ServeWebSocket(w, r, claims.Username)
}

func getTokenFromRequest(r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	return ""
}
