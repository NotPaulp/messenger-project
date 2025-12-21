package mlanalyze

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"messenger-project/internal/handlers/openrouter"
	"messenger-project/internal/models"
	messages "messenger-project/internal/repository/api-gateway"
	"messenger-project/pkg/config"
	"messenger-project/pkg/logger"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type MLAnalyzer struct {
	kafkaConsumer *kafka.Reader
	config        *config.Config
	logger        *logger.Logger
}

func NewMLAnalyzer(
	kafkaConsumer *kafka.Reader,
	cfg *config.Config,
	log *logger.Logger,
) *MLAnalyzer {
	return &MLAnalyzer{
		kafkaConsumer: kafkaConsumer,
		config:        cfg,
		logger:        log,
	}
}

func (a *MLAnalyzer) Start(ctx context.Context) {
	if a.logger == nil { // ✅ Safety check
		a.logger = logger.New(a.config.DebugMode)
	}
	a.logger.Info("Starting ML Analyzer service...")

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Context canceled, stopping ML Analyzer")
			return

		default:
			msg, err := a.kafkaConsumer.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					a.logger.Info("Context closed, stopping")
					return
				}
				a.logger.Error("Error reading message from Kafka: %v", err)
				continue
			}

			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				a.logger.Error("Error unmarshaling message: %v", err)
				continue
			}

			a.logger.Error("Processing message ID=%d from %s to %s",
				message.ID, message.SenderUsername, message.ReceiverUsername)

			analysis, err := a.analyzeMessage(message)
			if err != nil {
				a.logger.Error("Error analyzing message %d: %v", message.ID, err)
				continue
			}

			if err := a.publishAnalysisResult(analysis); err != nil {
				a.logger.Error("Error publishing analysis result: %v", err)
				continue
			}

			a.logger.Info("Analysis completed for message %d: category=%s, toxicity=%.2f",
				message.ID, analysis.Category, analysis.ToxicityScore)
		}
	}
}

func (a *MLAnalyzer) analyzeMessage(msg models.Message) (*models.AnalysisResult, error) {
	category, err := categorizeMessage(msg.Body)
	if err != nil {
		a.logger.Error("Categorization error for message %d: %v, using fallback", msg.ID, err)
		category = "general"
	}

	toxicityScore, err := evaluateToxicity(msg.Body)
	if err != nil {
		a.logger.Error("Toxicity evaluation error for message %d: %v, using fallback", msg.ID, err)
		toxicityScore = 0.0
	}

	isToxic := toxicityScore > 0.7

	result := &models.AnalysisResult{
		MessageID:      msg.ID,
		Category:       category,
		ToxicityScore:  toxicityScore,
		IsToxic:        isToxic,
		AnalyzedAt:     time.Now(),
		SenderUsername: msg.SenderUsername,
	}

	return result, nil
}

func (a *MLAnalyzer) publishAnalysisResult(
	result *models.AnalysisResult,
) error {

	updateData := map[string]any{
		"category":       result.Category,
		"toxicity_score": result.ToxicityScore,
		"toxic":          result.IsToxic,
		"analyzed_at":    result.AnalyzedAt,
	}
	if err := messages.UpdateMessageByID(result.MessageID, updateData); err != nil {
		return fmt.Errorf("error updating message in MongoDB: %w", err)
	}

	a.logger.Info("Published analysis result for message %d to Mongo", result.MessageID)
	return nil
}

const (
	CategorizationPrompt = `
You are a text classification model.

Task:
Given a short user message, choose exactly one category from this list:
question, statement, command, greeting, farewell, complaint, compliment, request, general, swear, spam, other.

Rules:
- Answer in English.
- Output MUST be valid JSON in this exact format:
{"category": "<one_of_the_categories_above>"}
- Do NOT add any other text before or after the JSON.

Message:
"%s"
`

	ToxicityPrompt = `
You are a toxicity scoring model.

Task:
Given a short user message, evaluate its toxicity on a scale from 0.0 (not toxic) to 1.0 (extremely toxic).

Rules:
- Output MUST be valid JSON in this exact format:
{"toxicity": <number_between_0.0_and_1.0>}
- Do NOT add any other text before or after the JSON.

Message:
"%s"
`
)

type categoryJSON struct {
	Category string `json:"category"`
}

func categorizeMessage(body string) (string, error) {
	prompt := fmt.Sprintf(CategorizationPrompt, body)
	resp, err := openrouter.CallOpenRouter("You are a classifier for chat messages.", prompt)
	if err != nil {
		return "", err
	}
	resp = strings.TrimSpace(resp)
	log.Printf("MLAnalyzerService: categorization response for message '%s': %s", body, resp)

	var cj categoryJSON
	if err := json.Unmarshal([]byte(resp), &cj); err != nil {
		return "", fmt.Errorf("MLAnalyzerService: JSON parse error in categorization: %w", err)
	}
	cat := strings.ToLower(strings.TrimSpace(cj.Category))
	if cat == "" {
		return "", fmt.Errorf("MLAnalyzerService: empty category in response")
	}
	return cat, nil
}

type toxicityJSON struct {
	Toxicity float32 `json:"toxicity"`
}

func evaluateToxicity(body string) (float32, error) {
	prompt := fmt.Sprintf(ToxicityPrompt, body)
	resp, err := openrouter.CallOpenRouter("You are a toxicity scoring model.", prompt)
	if err != nil {
		return 0.0, err
	}
	resp = strings.TrimSpace(resp)
	log.Printf("MLAnalyzerService: toxicity response for message '%s': %s", body, resp)

	var tj toxicityJSON
	if err := json.Unmarshal([]byte(resp), &tj); err != nil {
		return 0, fmt.Errorf("MLAnalyzerService: JSON parse error in toxicity: %w", err)
	}
	return tj.Toxicity, nil
}
