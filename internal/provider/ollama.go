package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndreaPallotta/civet/internal/config"
)

type OllamaProvider struct {
	endpoint string
	model    string
}

func NewOllamaProvider(cfg *config.Config) *OllamaProvider {
	endpoint := cfg.AI.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/generate"
	}
	model := cfg.AI.Model
	if model == "" {
		model = "llama3"
	}
	return &OllamaProvider{
		endpoint: endpoint,
		model:    model,
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

type ollamaRequest struct {
	Model  string `json:"model"`
	System string `json:"system"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (p *OllamaProvider) Analyze(ctx context.Context, req *AnalysisRequest) (*AnalysisResponse, error) {
	payload := ollamaRequest{
		Model:  p.model,
		System: systemPrompt,
		Prompt: buildContextPrompt(req),
		Stream: false,
		Format: "json",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status: %s", resp.Status)
	}

	var oResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return nil, err
	}

	if oResp.Response == "" {
		return nil, fmt.Errorf("empty response from Ollama")
	}

	var analysis AnalysisResponse
	if err := json.Unmarshal([]byte(oResp.Response), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from Ollama: %w\nRaw text: %s", err, oResp.Response)
	}

	return &analysis, nil
}
