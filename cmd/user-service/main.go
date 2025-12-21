package main

import (
	"context"
	"encoding/json"
	"log"
	"messenger-project/internal/database"
	grpcServerConstructor "messenger-project/internal/grpc"
	"messenger-project/internal/handlers/auth"
	"messenger-project/internal/redis"
	"messenger-project/pkg/config"
	"messenger-project/pkg/logger"
	pb "messenger-project/proto/user"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

	redis.Init()

	grpcLis, err := net.Listen("tcp", "0.0.0.0:"+cfg.GrpcPort)
	if err != nil {
		log.Error("gRPC listen: %v", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcServerConstructor.NewUserServer())

	go func() {
		log.Info("gRPC User Service on %s", "0.0.0.0:"+cfg.GrpcPort)
		grpcServer.Serve(grpcLis)
	}()

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/register", auth.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/api/logout", auth.LogoutHandler).Methods("POST")

	server := &http.Server{
		Addr:         ":" + cfg.UserServicePort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("User Service starting on http://localhost:%s", cfg.UserServicePort)
		log.Info("Health check: http://localhost:%s/health", cfg.UserServicePort)
		if err := server.ListenAndServe(); err != nil {
			log.Error("Server error:" + err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error:" + err.Error())
	}
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
		"service":   "user-service",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
