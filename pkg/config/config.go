package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort        string
	JWTSecret         string
	DebugMode         bool
	DATABASE_URL      string
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	REDIS_URL         string
	REDIS_PASSWORD    string
}

func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "default-secret-change-in-production"),
		DebugMode:         getEnvAsBool("DEBUG", true),
		DATABASE_URL:      getEnv("DATABASE_URL", "postgres://user:password@postgres:5432/db?sslmode=disable"),
		POSTGRES_DB:       getEnv("POSTGRES_DB", "db"),
		POSTGRES_USER:     getEnv("POSTGRES_USER", "user"),
		POSTGRES_PASSWORD: getEnv("POSTGRES_USER", "user"),
		REDIS_URL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		REDIS_PASSWORD:    getEnv("REDIS_PASSWORD", "password"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
