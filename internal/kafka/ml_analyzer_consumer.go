package kafka

import (
	"messenger-project/pkg/config"

	"github.com/segmentio/kafka-go"
)

type MLAnalyzerConsumer struct {
	reader *kafka.Reader
}

func NewMLAnalyzerConsumer(cfg *config.Config) *MLAnalyzerConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KAFKA_BROKERS},
		Topic:   cfg.KAFKA_TOPIC_ML_ANALYZE_MESSAGES,
		GroupID: "ml-analyzer-consumer-group",
	})
	return &MLAnalyzerConsumer{
		reader: reader,
	}
}

func (c *MLAnalyzerConsumer) Close() error {
	return c.reader.Close()
}
