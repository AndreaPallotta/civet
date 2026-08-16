package gitlab

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingEnvironmentRule{})
}

type MissingEnvironmentRule struct{}

func (r *MissingEnvironmentRule) ID() string { return "GL-007" }
func (r *MissingEnvironmentRule) Name() string { return "Missing Environment Declaration" }
func (r *MissingEnvironmentRule) Category() rules.Category { return rules.CategoryObservability }
func (r *MissingEnvironmentRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingEnvironmentRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *MissingEnvironmentRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		lowerName := strings.ToLower(name)
		if strings.Contains(lowerName, "deploy") || strings.Contains(lowerName, "release") || strings.Contains(lowerName, "publish") {
			if job.Environment == "" {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Deploy job is missing an 'environment' declaration.",
					Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
					Suggestion: "Declare 'environment' to track deployments and enable DORA metrics.",
				})
			}
		}
	}
	return findings
}
