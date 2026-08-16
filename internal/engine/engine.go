package engine

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/config"
	"github.com/AndreaPallotta/civet/internal/rules"
)

// Run is an alias for Analyze.
func Run(pipeline *rules.Pipeline, cfg *config.Config) *Report {
	return Analyze(pipeline, cfg)
}

// Analyze runs applicable rules against the pipeline and calculates scores.
func Analyze(pipeline *rules.Pipeline, cfg *config.Config) *Report {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	platformRules := rules.ForPlatform(pipeline.Platform)
	var activeRules []rules.Rule

	for _, r := range platformRules {
		if !cfg.IsDisabled(r.ID()) {
			activeRules = append(activeRules, r)
		}
	}

	var findings []rules.Finding
	for _, r := range activeRules {
		ruleFindings := r.Check(pipeline)
		
		// Apply severity overrides if any
		if cfg.Rules.SeverityOverrides != nil {
			if overrideStr, ok := cfg.Rules.SeverityOverrides[r.ID()]; ok {
				overrideStr = strings.ToLower(overrideStr)
				var newSeverity rules.Severity
				overrideValid := true
				
				switch overrideStr {
				case "info":
					newSeverity = rules.SeverityInfo
				case "warning":
					newSeverity = rules.SeverityWarning
				case "error":
					newSeverity = rules.SeverityError
				default:
					overrideValid = false
				}
				
				if overrideValid {
					for i := range ruleFindings {
						ruleFindings[i].Severity = newSeverity
					}
				}
			}
		}
		
		findings = append(findings, ruleFindings...)
	}

	categoryScores := make(map[rules.Category]int)
	for _, cat := range rules.AllCategories() {
		categoryScores[cat] = 100
	}

	for _, f := range findings {
		penalty := 0
		switch f.Severity {
		case rules.SeverityError:
			penalty = 15
		case rules.SeverityWarning:
			penalty = 8
		case rules.SeverityInfo:
			penalty = 3
		}
		
		categoryScores[f.Category] -= penalty
		if categoryScores[f.Category] < 0 {
			categoryScores[f.Category] = 0
		}
	}

	totalScore := 0
	count := 0
	
	stringCategoryScores := make(map[string]int)
	for cat, score := range categoryScores {
		totalScore += score
		stringCategoryScores[cat.String()] = score
		count++
	}

	overallScore := 0
	if count > 0 {
		overallScore = totalScore / count
	}

	return &Report{
		FilePath:       pipeline.FilePath,
		Platform:       pipeline.Platform,
		OverallScore:   overallScore,
		CategoryScores: stringCategoryScores,
		Findings:       findings,
		RawYAML:        pipeline.Raw,
	}
}
