package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"
const deepseekModel = "nex-agi/deepseek-v3.1-nex-n1:free"

type ORMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ORChatRequest struct {
	Model    string      `json:"model"`
	Messages []ORMessage `json:"messages"`
}

type ORChatChoice struct {
	Message ORMessage `json:"message"`
}

type ORChatResponse struct {
	Choices []ORChatChoice `json:"choices"`
}

var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

func CallOpenRouter(systemPrompt, userPrompt string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY is not set")
	}

	reqBody := ORChatRequest{
		Model: deepseekModel,
		Messages: []ORMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openRouterURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "https://your-app.local")
	req.Header.Set("X-Title", "Messenger ML Analyzer")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openrouter status %d: %s", resp.StatusCode, string(raw))
	}

	var orResp ORChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return "", err
	}

	if len(orResp.Choices) == 0 {
		return "", fmt.Errorf("openrouter: empty choices")
	}

	return orResp.Choices[0].Message.Content, nil
}
