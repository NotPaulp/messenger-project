package auth

import (
	"fmt"
	"messenger-project/internal/redis"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID   string
	Username string
	jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	return signedToken, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	isBlacklisted, err := redis.IsJWTBlacklisted(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	if isBlacklisted {
		return nil, fmt.Errorf("token is blacklisted")
	}

	return claims, nil
}
