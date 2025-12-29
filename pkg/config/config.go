package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort                     string
	GrpcPort                       string
	UserGrpcAddr                   string
	MLAnalyzerPort                 string
	UserServicePort                string
	JWTSecret                      string
	DebugMode                      bool
	DATABASE_URL                   string
	POSTGRES_DB                    string
	POSTGRES_USER                  string
	POSTGRES_PASSWORD              string
	REDIS_URL                      string
	REDIS_PASSWORD                 string
	MONGO_URL                      string
	MONGO_DB                       string
	MONGO_MESSAGES_COLLECTION      string
	MONGO_POSTS_COLLECTION         string
	KAFKA_BROKERS                  string
	KAFKA_TOPIC_CONTENT            string
	KAFKA_TOPIC_ML_ANALYZE_CONTENT string
	WEB_SOCKET_READ_BUFFER_SIZE    int
	WEB_SOCKET_WRITE_BUFFER_SIZE   int
}

func Load() *Config {
	return &Config{
		ServerPort:                     getEnv("SERVER_PORT", "8080"),
		MLAnalyzerPort:                 getEnv("ML_ANALYZER_PORT", "8081"),
		UserServicePort:                getEnv("USER_SERVICE_PORT", "8082"),
		GrpcPort:                       getEnv("GRPC_PORT", "8083"),
		UserGrpcAddr:                   getEnv("USER_GRPC_ADDR", "user-service:8083"),
		JWTSecret:                      getEnv("JWT_SECRET", "default-secret-change-in-production"),
		DebugMode:                      getEnvAsBool("DEBUG", true),
		DATABASE_URL:                   getEnv("DATABASE_URL", "postgres://user:password@postgres:5432/db?sslmode=disable"),
		POSTGRES_DB:                    getEnv("POSTGRES_DB", "db"),
		POSTGRES_USER:                  getEnv("POSTGRES_USER", "user"),
		POSTGRES_PASSWORD:              getEnv("POSTGRES_USER", "user"),
		REDIS_URL:                      getEnv("REDIS_URL", "redis://localhost:6379"),
		REDIS_PASSWORD:                 getEnv("REDIS_PASSWORD", "password"),
		MONGO_URL:                      getEnv("MONGO_URL", "mongodb://mongo:27017"),
		MONGO_DB:                       getEnv("MONGO_DB", "messenger"),
		MONGO_MESSAGES_COLLECTION:      getEnv("MONGO_MESSAGES_COLLECTION", "messages"),
		MONGO_POSTS_COLLECTION:         getEnv("MONGO_POSTS_COLLECTION", "posts"),
		KAFKA_BROKERS:                  getEnv("KAFKA_BROKERS", "localhost:9092"),
		KAFKA_TOPIC_CONTENT:            getEnv("KAFKA_TOPIC_CONTENT", "content"),
		KAFKA_TOPIC_ML_ANALYZE_CONTENT: getEnv("KAFKA_TOPIC_ML_ANALYZE_CONTENT", "ml_analyze_content"),
		WEB_SOCKET_READ_BUFFER_SIZE:    1024,
		WEB_SOCKET_WRITE_BUFFER_SIZE:   1024,
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

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
