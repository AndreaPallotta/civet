package github

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingCacheActionRule{})
}

type MissingCacheActionRule struct{}

func (r *MissingCacheActionRule) ID() string { return "GH-008" }
func (r *MissingCacheActionRule) Name() string { return "Missing Cache Action" }
func (r *MissingCacheActionRule) Category() rules.Category { return rules.CategoryPerformance }
func (r *MissingCacheActionRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingCacheActionRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *MissingCacheActionRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	hasCache := false
	hasDependencies := false

	for _, job := range pipeline.Jobs {
		for _, step := range job.Steps {
			if strings.HasPrefix(step.Uses, "actions/cache") {
				hasCache = true
			}
			if strings.HasPrefix(step.Uses, "actions/setup-") {
				if _, ok := step.With["cache"]; ok {
					hasCache = true
				}
			}

			runLower := strings.ToLower(step.Run)
			if strings.Contains(runLower, "npm ") || strings.Contains(runLower, "yarn ") || strings.Contains(runLower, "pip ") || strings.Contains(runLower, "go ") {
				hasDependencies = true
			}
		}
	}

	if hasDependencies && !hasCache {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline appears to install dependencies but does not use caching.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Use actions/cache or actions/setup-* with the 'cache' option to speed up workflow execution.",
			},
		}
	}

	return nil
}
