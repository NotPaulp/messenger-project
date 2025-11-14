package redis

import (
	"context"
	"fmt"
	"messenger-project/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type Claims struct {
	UserID   string
	Username string
	jwt.RegisteredClaims
}

var Client *redis.Client
var ctx = context.Background() // smth like OS.ENV + expiration time

func Init() error {
	cfg := config.Load()
	redisURL := cfg.REDIS_URL

	opts, err := redis.ParseURL(redisURL) // opts - struct with URL vals
	if err != nil {
		return fmt.Errorf("Failed to parse Redis URL: %v", err)
	}

	Client = redis.NewClient(opts)

	_, err = Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}

func BlacklistJWT(token string) error {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := &Claims{}
	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return fmt.Errorf("Failed to parse JWT token: %v", err)
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl < 0 {
		return fmt.Errorf("JWT token is already expired")
	}

	err = Client.Set(ctx, token, "blacklisted", ttl).Err()
	if err != nil {
		return fmt.Errorf("Failed to blacklist JWT token: %v", err)
	}

	return nil
}

func IsJWTBlacklisted(token string) (bool, error) {
	val, err := Client.Get(ctx, token).Result()
	if err != redis.Nil && err != nil {
		return false, fmt.Errorf("Error while getting the token out of redis db %v", err)
	}

	switch val {
	case "":
		return false, nil
	case "blacklisted":
		return true, nil
	default:
		return false, fmt.Errorf("Failed to check if JWT token is blacklisted, val: %v", val)
	}
}
