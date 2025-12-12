package mlanalyze

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"messenger-project/internal/handlers/openrouter"
	"messenger-project/internal/repository"
	"strings"
	"time"
)

func Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("MLAnalyzerService: context canceled, stopping")
			return

		case <-ticker.C:
			analyze()
		}
	}
}

func analyze() {
	where := map[string]any{
		"category": "",
	}
	messages, err := repository.GetAllMessagesWhere(where)
	if err != nil {
		log.Printf("MLAnalyzerService: error getting messages: %v", err)
		return
	}

	if len(messages) == 0 {
		return
	}

	log.Printf("MLAnalyzerService: found %d messages to process", len(messages))

	for _, message := range messages {
		category, err := categorizeMessage(message.Body)

		if err != nil {
			log.Printf("MLAnalyzerService: error categorizing message %d: %v", message.ID, err)
			continue
		}

		toxicityScore, err := evaluateToxicity(message.Body)
		if err != nil {
			log.Printf("MLAnalyzerService: error evaluating toxicity for message %d: %v", message.ID, err)
			continue
		}

		toUpdate := map[string]any{
			"category":       category,
			"toxicity_score": toxicityScore,
			"toxic":          toxicityScore > 0.7,
		}

		if err := repository.UpdateMessageStatus(message.ID, toUpdate); err != nil {
			log.Printf("MLAnalyzerService: error setting category and toxicity for message %d: %v", message.ID, err)
		}
	}
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
