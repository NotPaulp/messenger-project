package auth

import (
	"context"
	"messenger-project/internal/redis"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Auth header is missing", http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		isBlacklisted, err := redis.IsJWTBlacklisted(tokenString)
		if err != nil {
			http.Error(w, "Internal error"+err.Error(), http.StatusInternalServerError)
			return
		}

		if isBlacklisted {
			http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
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
