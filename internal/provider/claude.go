package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndreaPallotta/civet/internal/config"
)

type ClaudeProvider struct {
	model  string
	apiKey string
}

func NewClaudeProvider(cfg *config.Config) *ClaudeProvider {
	model := cfg.AI.Model
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &ClaudeProvider{
		model:  model,
		apiKey: getEnvOrDefault(cfg.AI.APIKeyEnv, "ANTHROPIC_API_KEY"),
	}
}

func (p *ClaudeProvider) Name() string {
	return "claude"
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	System    string          `json:"system"`
	Messages  []claudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (p *ClaudeProvider) Analyze(ctx context.Context, req *AnalysisRequest) (*AnalysisResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key is required for Claude")
	}

	payload := claudeRequest{
		Model:     p.model,
		System:    systemPrompt,
		MaxTokens: 2048,
		Messages: []claudeMessage{
			{Role: "user", Content: buildContextPrompt(req)},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("claude API returned status: %s", resp.Status)
	}

	var cResp claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return nil, err
	}

	if len(cResp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	text := cResp.Content[0].Text

	var analysis AnalysisResponse
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from Claude: %w\nRaw text: %s", err, text)
	}

	return &analysis, nil
}
