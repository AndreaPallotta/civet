package universal

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingJobTimeoutRule{})
}

type MissingJobTimeoutRule struct{}

func (r *MissingJobTimeoutRule) ID() string { return "UNI-003" }
func (r *MissingJobTimeoutRule) Name() string { return "Missing Job Timeout" }
func (r *MissingJobTimeoutRule) Category() rules.Category { return rules.CategoryReliability }
func (r *MissingJobTimeoutRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingJobTimeoutRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}
func (r *MissingJobTimeoutRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	for _, job := range pipeline.Jobs {
		if job.Timeout == "" {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job is missing a timeout configuration",
				Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
				Suggestion: "Set a timeout (GitLab timeout: or GitHub timeout-minutes:) to prevent jobs from hanging indefinitely.",
			})
		}
	}
	return findings
}
