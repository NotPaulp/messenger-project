package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"messenger-project/internal/database"
	mlanalyze "messenger-project/internal/ml-analyze"
	"messenger-project/pkg/config"
	"messenger-project/pkg/logger"

	"github.com/segmentio/kafka-go"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.DebugMode)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := database.InitMongo(ctx); err != nil {
		log.Error("Mongo Error: " + err.Error())
		os.Exit(1)
	}
	defer database.DisconnectMongo(ctx)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KAFKA_BROKERS},
		Topic:   cfg.KAFKA_TOPIC_ML_ANALYZE_MESSAGES,
		GroupID: "ml-analyzer-consumer-group",
	})
	defer reader.Close()
	analyzer := mlanalyze.NewMLAnalyzer(reader, cfg, log)
	go analyzer.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("ML Analyzer running...")
	<-sigChan

	log.Info("Shutting down...")
	cancel()
}
