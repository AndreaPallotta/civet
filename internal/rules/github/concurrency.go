package github

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&NoConcurrencyGroupRule{})
}

type NoConcurrencyGroupRule struct{}

func (r *NoConcurrencyGroupRule) ID() string { return "GH-003" }
func (r *NoConcurrencyGroupRule) Name() string { return "No Concurrency Group" }
func (r *NoConcurrencyGroupRule) Category() rules.Category { return rules.CategoryReliability }
func (r *NoConcurrencyGroupRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *NoConcurrencyGroupRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *NoConcurrencyGroupRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	if pipeline.Concurrency == nil || pipeline.Concurrency.Group == "" {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline does not define a concurrency group.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Use 'concurrency' to prevent duplicate pipeline runs for the same ref.",
			},
		}
	}

	return nil
}
