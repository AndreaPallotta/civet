package github

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&EnvSecretsExposureRule{})
}

type EnvSecretsExposureRule struct{}

func (r *EnvSecretsExposureRule) ID() string { return "GH-007" }
func (r *EnvSecretsExposureRule) Name() string { return "Workflow-Level Secret Exposure" }
func (r *EnvSecretsExposureRule) Category() rules.Category { return rules.CategorySecurity }
func (r *EnvSecretsExposureRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *EnvSecretsExposureRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *EnvSecretsExposureRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		for key, val := range job.Variables {
			if strings.Contains(val, "${{ secrets.") {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Secret exposed at the job or workflow environment level.",
					Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
					Suggestion: "Map secrets to the environment of the specific steps that need them to reduce exposure risk. Avoid global variables key: " + key,
				})
			}
		}
	}
	return findings
}
