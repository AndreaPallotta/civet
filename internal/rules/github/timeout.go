package github

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingTimeoutRule{})
}

type MissingTimeoutRule struct{}

func (r *MissingTimeoutRule) ID() string { return "GH-005" }
func (r *MissingTimeoutRule) Name() string { return "Missing Timeout Minutes" }
func (r *MissingTimeoutRule) Category() rules.Category { return rules.CategoryReliability }
func (r *MissingTimeoutRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingTimeoutRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *MissingTimeoutRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		if job.Timeout == "" {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job is missing 'timeout-minutes'.",
				Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
				Suggestion: "Set a specific timeout to avoid runaway jobs using the excessive 360-minute default.",
			})
		}
	}
	return findings
}
