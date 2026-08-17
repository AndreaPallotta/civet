package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/AndreaPallotta/civet/internal/engine"
	"github.com/AndreaPallotta/civet/internal/rules"
)

// LLMReporter formats pipeline audits into a dense, token-optimized format
// specifically designed for LLM agent context windows and reasoning engines.
type LLMReporter struct{}

func NewLLMReporter() *LLMReporter {
	return &LLMReporter{}
}

// LLMContext represents the structured context block injected into LLM prompts.
type LLMContext struct {
	Header      string                 `json:"system_prompt_context"`
	Pipeline    LLMPipelineInfo        `json:"pipeline"`
	Scores      map[string]int         `json:"category_scores"`
	TotalScore  int                    `json:"overall_score"`
	HealthGrade string                 `json:"health_grade"`
	Findings    []LLMFinding           `json:"findings"`
	ActionItems []string               `json:"recommended_llm_actions"`
}

type LLMPipelineInfo struct {
	FilePath string `json:"file_path"`
	Platform string `json:"platform"`
}

type LLMFinding struct {
	RuleID      string         `json:"rule_id"`
	RuleName    string         `json:"rule_name"`
	Severity    rules.Severity `json:"severity"`
	Category    rules.Category `json:"category"`
	Location    rules.Location `json:"location"`
	Description string         `json:"issue_description"`
	ProposedFix string         `json:"suggested_fix"`
}

func (r *LLMReporter) Write(w io.Writer, rep *engine.Report) error {
	grade := "HEALTHY"
	if rep.OverallScore < 60 {
		grade = "CRITICAL_ATTENTION_REQUIRED"
	} else if rep.OverallScore < 85 {
		grade = "NEEDS_IMPROVEMENT"
	}

	var findings []LLMFinding
	var actionItems []string

	for _, f := range rep.Findings {
		findings = append(findings, LLMFinding{
			RuleID:      f.RuleID,
			RuleName:    f.RuleName,
			Severity:    f.Severity,
			Category:    f.Category,
			Location:    f.Location,
			Description: f.Message,
			ProposedFix: f.Suggestion,
		})

		action := fmt.Sprintf("[%s] Resolve %s: %s", f.Severity.String(), f.RuleID, f.Suggestion)
		actionItems = append(actionItems, action)
	}

	contextData := LLMContext{
		Header: "Civet Deterministic CI/CD Quality Audit. Use the structured findings below to review the user's pipeline and propose direct code edits.",
		Pipeline: LLMPipelineInfo{
			FilePath: rep.FilePath,
			Platform: rep.Platform.String(),
		},
		Scores:      rep.CategoryScores,
		TotalScore:  rep.OverallScore,
		HealthGrade: grade,
		Findings:    findings,
		ActionItems: actionItems,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(contextData)
}
