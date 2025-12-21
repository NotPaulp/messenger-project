package llama

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type LlamaRequest struct {
	Prompt      string   `json:"prompt"`
	NPredict    int      `json:"n_predict"`
	Temperature float32  `json:"temperature"`
	Stop        []string `json:"stop,omitempty"`
}

type LlamaResponse struct {
	Content string `json:"content"`
}

func CallLlama(prompt string) (string, error) {
	body := LlamaRequest{
		Prompt:      prompt,
		NPredict:    128,
		Temperature: 0.1,
		Stop:        []string{"<|eot_id|>"},
	}

	b, _ := json.Marshal(body)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://192.168.1.2:8088/completion",
		bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out LlamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Content, nil
}
