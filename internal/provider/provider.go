package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/AndreaPallotta/civet/internal/config"
	"github.com/AndreaPallotta/civet/internal/rules"
)

type AnalysisRequest struct {
	PipelineYAML string
	Platform     string
	Findings     []rules.Finding
}

type AnalysisResponse struct {
	Summary         string   `json:"summary"`
	Recommendations []string `json:"recommendations"`
	RiskLevel       string   `json:"risk_level"`
}

type Provider interface {
	Name() string
	Analyze(ctx context.Context, req *AnalysisRequest) (*AnalysisResponse, error)
}

func NewProvider(cfg *config.Config) (Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is missing")
	}

	switch cfg.AI.Provider {
	case "claude":
		return NewClaudeProvider(cfg), nil
	case "openai":
		return NewOpenAIProvider(cfg), nil
	case "gemini":
		return NewGeminiProvider(cfg), nil
	case "ollama":
		return NewOllamaProvider(cfg), nil
	case "":
		return nil, fmt.Errorf("AI provider not configured")
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.AI.Provider)
	}
}

func getEnvOrDefault(envKey, defaultVal string) string {
	if envKey == "" {
		envKey = defaultVal
	}
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	if val := os.Getenv(defaultVal); val != "" {
		return val
	}
	return ""
}

const systemPrompt = "You are Civet, a CI/CD pipeline analysis assistant. You have been given a pipeline configuration file and the results of a static rules analysis. Provide a concise summary of the pipeline's health, specific actionable recommendations beyond what the rules engine already found, and an overall risk level (low, medium, high, or critical). Respond in JSON format with keys: summary, recommendations (array of strings), risk_level."

func buildContextPrompt(req *AnalysisRequest) string {
	findingsStr := "Findings:\n"
	for _, f := range req.Findings {
		findingsStr += fmt.Sprintf("- [%s] %s: %s\n", f.Severity.String(), f.RuleID, f.Message)
	}

	return fmt.Sprintf("Platform: %s\n\nPipeline YAML:\n%s\n\n%s", req.Platform, req.PipelineYAML, findingsStr)
}
