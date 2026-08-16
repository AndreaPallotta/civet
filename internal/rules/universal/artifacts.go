package universal

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingArtifactExpirationRule{})
}

type MissingArtifactExpirationRule struct{}

func (r *MissingArtifactExpirationRule) ID() string { return "UNI-006" }
func (r *MissingArtifactExpirationRule) Name() string { return "Missing Artifact Expiration" }
func (r *MissingArtifactExpirationRule) Category() rules.Category { return rules.CategoryPerformance }
func (r *MissingArtifactExpirationRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingArtifactExpirationRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}

func (r *MissingArtifactExpirationRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	for _, job := range pipeline.Jobs {
		if job.Artifacts != nil && job.Artifacts.ExpireIn == "" {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job produces artifacts but does not configure an expiration time",
				Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
				Suggestion: "Set an expiration for artifacts to prevent unbounded storage growth.",
			})
		}
	}
	return findings
}
