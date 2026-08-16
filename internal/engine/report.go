package engine

import "github.com/AndreaPallotta/civet/internal/rules"

type AIAnalysis struct {
	Summary         string   `json:"summary"`
	Recommendations []string `json:"recommendations"`
	RiskLevel       string   `json:"risk_level"`
}

type Report struct {
	FilePath       string             `json:"file"`
	Platform       rules.Platform     `json:"platform"`
	OverallScore   int                `json:"score"`
	CategoryScores map[string]int     `json:"categories"`
	Findings       []rules.Finding    `json:"findings"`
	RawYAML        string             `json:"-"`
	AIAnalysis     *AIAnalysis        `json:"ai_analysis,omitempty"`
}
