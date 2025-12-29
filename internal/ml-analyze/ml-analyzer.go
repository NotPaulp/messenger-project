package mlanalyze

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"messenger-project/internal/handlers/openrouter"
	"messenger-project/internal/models"
	api_gateway "messenger-project/internal/repository/api-gateway"
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
	if a.logger == nil {
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
			contentType := ""
			for _, h := range msg.Headers {
				if h.Key == "content_type" {
					contentType = string(h.Value)
					break
				}
			}

			switch contentType {
			case "message":
				a.handleMessageAnalysis(msg.Value)
			case "post":
				a.handlePostAnalysis(msg.Value)
			case "comment":
				a.handleCommentAnalysis(msg.Value)
			default:
				a.logger.Error("Unknown content type: %s", contentType)
			}

		}
	}
}

func (a *MLAnalyzer) handleMessageAnalysis(value []byte) {
	var message models.Message
	if err := json.Unmarshal(value, &message); err != nil {
		a.logger.Error("Error unmarshaling message: %v", err)
		return
	}

	a.logger.Info("Processing message ID=%d from %s to %s",
		message.ID, message.SenderUsername, message.ReceiverUsername)

	msgAnalysis, err := a.analyzeMessage(message)
	if err != nil {
		a.logger.Error("Error analyzing message %d: %v", message.ID, err)
		return
	}

	if err := a.publishMessageAnalysisResult(msgAnalysis); err != nil {
		a.logger.Error("Error publishing message analysis result: %v", err)
		return
	}

	a.logger.Info("Analysis completed for message %d: category=%s, toxicity=%.2f",
		message.ID, msgAnalysis.Category, msgAnalysis.ToxicityScore)
}

func (a *MLAnalyzer) handlePostAnalysis(value []byte) {
	var post models.Post
	if err := json.Unmarshal(value, &post); err != nil {
		a.logger.Error("Error unmarshaling post: %v", err)
		return
	}

	a.logger.Info("Processing post ID=%d from %s",
		post.ID, post.AuthorUsername)

	postAnalysis, err := a.analyzePost(post)
	if err != nil {
		a.logger.Error("Error analyzing post %d: %v", post.ID, err)
		return
	}

	if err := a.publishPostAnalysisResult(postAnalysis); err != nil {
		a.logger.Error("Error publishing post analysis result: %v", err)
		return
	}

	a.logger.Info("Analysis completed for post %d: categories=%v",
		post.ID, postAnalysis.Categories)
}

func (a *MLAnalyzer) handleCommentAnalysis(value []byte) {
	var comment models.Comment
	if err := json.Unmarshal(value, &comment); err != nil {
		a.logger.Error("Error unmarshaling comment: %v", err)
		return
	}

	a.logger.Info("Processing comment ID=%d on post %d from %s",
		comment.ID, comment.PostID, comment.AuthorUsername)

	commentAnalysis, err := a.analyzeComment(comment)
	if err != nil {
		a.logger.Error("Error analyzing comment %d: %v", comment.ID, err)
		return
	}

	if err := a.publishCommentAnalysisResult(commentAnalysis); err != nil {
		a.logger.Error("Error publishing comment analysis result: %v", err)
		return
	}

	a.logger.Info("Analysis completed for comment %d: spam=%v, toxicity=%.2f",
		comment.ID, commentAnalysis.Spam, commentAnalysis.ToxicityScore)
}

func (a *MLAnalyzer) analyzeMessage(msg models.Message) (*models.MessageAnalysisResult, error) {
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

	result := &models.MessageAnalysisResult{
		MessageID:      msg.ID,
		Category:       category,
		ToxicityScore:  toxicityScore,
		IsToxic:        isToxic,
		AnalyzedAt:     time.Now(),
		SenderUsername: msg.SenderUsername,
	}

	return result, nil
}

func (a *MLAnalyzer) analyzePost(post models.Post) (*models.PostAnalysisResult, error) {
	categories, err := categorizePost(post.Body)
	if err != nil {
		a.logger.Error("Categorization error for post %d: %v, using fallback", post.ID, err)
		categories = []string{"general"}
	}
	result := &models.PostAnalysisResult{
		PostID:         post.ID,
		Categories:     categories,
		AnalyzedAt:     time.Now(),
		AuthorUsername: post.AuthorUsername,
	}
	return result, nil
}

func (a *MLAnalyzer) analyzeComment(comment models.Comment) (*models.CommentAnalysisResult, error) {
	spam, err := detectSpamInComment(comment.Body)
	if err != nil {
		a.logger.Error("Spam detection error for comment %d: %v, using fallback", comment.ID, err)
		spam = false
	}

	comment.ToxicityScore, err = evaluateToxicity(comment.Body)
	if err != nil {
		a.logger.Error("Toxicity evaluation error for comment %d: %v, using fallback", comment.ID, err)
		comment.ToxicityScore = 0.0
	}
	isToxic := comment.ToxicityScore > 0.7

	result := &models.CommentAnalysisResult{
		CommentID:      comment.ID,
		PostID:         comment.PostID,
		ToxicityScore:  comment.ToxicityScore,
		IsToxic:        isToxic,
		Spam:           spam,
		AnalyzedAt:     time.Now(),
		AuthorUsername: comment.AuthorUsername,
	}
	return result, nil
}

func (a *MLAnalyzer) publishMessageAnalysisResult(
	result *models.MessageAnalysisResult,
) error {

	updateData := map[string]any{
		"category":       result.Category,
		"toxicity_score": result.ToxicityScore,
		"toxic":          result.IsToxic,
		"analyzed_at":    result.AnalyzedAt,
	}
	if err := api_gateway.UpdateMessageByID(result.MessageID, updateData); err != nil {
		return fmt.Errorf("error updating message in MongoDB: %w", err)
	}

	a.logger.Info("Published analysis result for message %d to Mongo", result.MessageID)
	return nil
}

func (a *MLAnalyzer) publishPostAnalysisResult(
	result *models.PostAnalysisResult,
) error {

	updateData := map[string]any{
		"categories":  result.Categories,
		"analyzed_at": result.AnalyzedAt,
	}
	if err := api_gateway.UpdatePostByID(result.PostID, updateData); err != nil {
		return fmt.Errorf("error updating post in MongoDB: %w", err)
	}
	a.logger.Info("Published analysis result for post %d to Mongo", result.PostID)
	return nil
}

func (a *MLAnalyzer) publishCommentAnalysisResult(
	result *models.CommentAnalysisResult,
) error {

	updateData := map[string]any{
		"spam":           result.Spam,
		"toxicity_score": result.ToxicityScore,
		"toxic":          result.IsToxic,
		"analyzed_at":    result.AnalyzedAt,
	}
	if err := api_gateway.UpdateCommentByID(result.PostID, result.CommentID, updateData); err != nil {
		return fmt.Errorf("error updating comment in MongoDB: %w", err)
	}
	a.logger.Info("Published analysis result for comment %d to Mongo", result.CommentID)
	return nil
}

const (
	MessageCategorizationPrompt = `
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
Given a short user message/comment, evaluate its toxicity on a scale from 0.0 (not toxic) to 1.0 (extremely toxic).

Rules:
- Output MUST be valid JSON in this exact format:
{"toxicity": <number_between_0.0_and_1.0>}
- Do NOT add any other text before or after the JSON.

Message:
"%s"
`

	PostCategorizationPrompt = `
You are a text classification model.

Task:
Given a blog post content, assign it to one or more categories from this list:
technology, health, lifestyle, finance, travel, food, education, entertainment, sports, politics, science, fashion, general.

Rules:
- Answer in English.
- Output MUST be valid JSON in this exact format:
{"categories": ["<category1>", "<category2>", ...]}
- Do NOT add any other text before or after the JSON.

Post Content:
"%s"
`

	CommentSpamDetectionPrompt = `
You are a spam detection model.

Task:
Given a user comment, determine if it is spam or not.

Rules:
- Output MUST be valid JSON in this exact format:
{"spam": true} or {"spam": false}
- Do NOT add any other text before or after the JSON.

Comment:
"%s"
`
)

type categoryJSON struct {
	Category string `json:"category"`
}

func categorizeMessage(body string) (string, error) {
	prompt := fmt.Sprintf(MessageCategorizationPrompt, body)
	resp, err := openrouter.CallOpenRouter("You are a classifier for chat api_gateway", prompt)
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

type categoriesJSON struct {
	Categories []string `json:"categories"`
}

func categorizePost(body string) ([]string, error) {
	prompt := fmt.Sprintf(PostCategorizationPrompt, body)
	resp, err := openrouter.CallOpenRouter("You are a classifier for blog posts.", prompt)
	if err != nil {
		return nil, err
	}
	resp = strings.TrimSpace(resp)
	log.Printf("MLAnalyzerService: categorization response for post '%s': %s", body, resp)
	var cj categoriesJSON
	if err := json.Unmarshal([]byte(resp), &cj); err != nil {
		return nil, fmt.Errorf("MLAnalyzerService: JSON parse error in post categorization: %w", err)
	}
	if len(cj.Categories) == 0 {
		return nil, fmt.Errorf("MLAnalyzerService: empty categories in response")
	}
	return cj.Categories, nil
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

type spamJSON struct {
	Spam bool `json:"spam"`
}

func detectSpamInComment(body string) (bool, error) {
	prompt := fmt.Sprintf(CommentSpamDetectionPrompt, body)
	resp, err := openrouter.CallOpenRouter("You are a spam detection model.", prompt)
	if err != nil {
		return false, err
	}
	resp = strings.TrimSpace(resp)
	log.Printf("MLAnalyzerService: spam detection response for comment '%s': %s", body, resp)

	var sj spamJSON
	if err := json.Unmarshal([]byte(resp), &sj); err != nil {
		return false, fmt.Errorf("MLAnalyzerService: JSON parse error in spam detection: %w", err)
	}
	return sj.Spam, nil
}
