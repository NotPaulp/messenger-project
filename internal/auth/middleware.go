package auth

import (
	"context"
	"net/http"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Auth header is missing", http.StatusUnauthorized)
			return
		}
		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getTokenFromHeader(r *http.Request) string {
	authString := r.Header.Get("Authorization")
	if authString == "" {
		return ""
	}
	parts := strings.Fields(authString)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
