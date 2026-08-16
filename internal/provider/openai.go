package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndreaPallotta/civet/internal/config"
)

type OpenAIProvider struct {
	model  string
	apiKey string
}

func NewOpenAIProvider(cfg *config.Config) *OpenAIProvider {
	model := cfg.AI.Model
	if model == "" {
		model = "gpt-4o"
	}
	return &OpenAIProvider{
		model:  model,
		apiKey: getEnvOrDefault(cfg.AI.APIKeyEnv, "OPENAI_API_KEY"),
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model          string          `json:"model"`
	Messages       []openAIMessage `json:"messages"`
	ResponseFormat map[string]string `json:"response_format"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *OpenAIProvider) Analyze(ctx context.Context, req *AnalysisRequest) (*AnalysisResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key is required for OpenAI")
	}

	payload := openAIRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: buildContextPrompt(req)},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai API returned status: %s", resp.Status)
	}

	var oResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return nil, err
	}

	if len(oResp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	text := oResp.Choices[0].Message.Content

	var analysis AnalysisResponse
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from OpenAI: %w\nRaw text: %s", err, text)
	}

	return &analysis, nil
}
