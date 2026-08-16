package universal

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingRetryConfigRule{})
}

type MissingRetryConfigRule struct{}

func (r *MissingRetryConfigRule) ID() string { return "UNI-004" }
func (r *MissingRetryConfigRule) Name() string { return "Missing Retry Configuration" }
func (r *MissingRetryConfigRule) Category() rules.Category { return rules.CategoryReliability }
func (r *MissingRetryConfigRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *MissingRetryConfigRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}
func (r *MissingRetryConfigRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	for _, job := range pipeline.Jobs {
		if job.Retry == nil {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job is missing retry configuration for transient failures",
				Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
				Suggestion: "Configure retry for the job to handle temporary network or service failures.",
			})
		}
	}
	return findings
}
