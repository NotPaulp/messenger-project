package main

import (
	"context"
	"encoding/json"
	"log"
	middleware "messenger-project/internal/auth"
	"messenger-project/internal/database"
	"messenger-project/internal/handlers/auth"
	"messenger-project/internal/handlers/comments"
	"messenger-project/internal/handlers/messages"
	"messenger-project/internal/handlers/posts"
	"messenger-project/internal/redis"
	"messenger-project/pkg/config"
	"messenger-project/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.DebugMode)

	db, err := database.InitPostgres()
	if err != nil {
		log.Error("DB Error:" + err.Error())
		os.Exit(1)
	} else {
		defer db.Close()
	}

	database.CreateTableUsers(db)
	// database.CreateTableMessages(db)
	redis.Init()

	ctx := context.Background()
	if err := database.InitMongo(ctx); err != nil {
		log.Error("Mongo Error: " + err.Error())
		os.Exit(1)
	}
	defer database.DisconnectMongo(ctx)

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/register", auth.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/api/logout", auth.LogoutHandler).Methods("POST")
	router.Handle("/api/messages", middleware.Middleware(http.HandlerFunc(messages.SendHandler))).Methods("POST")
	router.Handle("/api/messages", middleware.Middleware(http.HandlerFunc(messages.GetHandler))).Methods("GET")
	router.Handle("/api/posts", middleware.Middleware(http.HandlerFunc(posts.PublishPostHandler))).Methods("POST")
	router.Handle("/api/posts", middleware.Middleware(http.HandlerFunc(posts.GetPostsHandler))).Methods("GET")
	router.Handle("/api/posts/comments", middleware.Middleware(http.HandlerFunc(comments.PublishCommentHandler))).Methods("POST")
	router.Handle("/api/posts/comments", middleware.Middleware(http.HandlerFunc(comments.GetCommentsHandler))).Methods("GET")

	//router.HandleFunc("/api/profile", profileHandler).Methods("GET")

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("API Gateway starting on http://localhost:%s", cfg.ServerPort)
		log.Info("Health check: http://localhost:%s/health", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil {
			log.Error("Server error:" + err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Info("Server stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]any{
		"status":    "healthy",
		"service":   "api-gateway",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
