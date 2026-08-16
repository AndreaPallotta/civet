package github

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&DeprecatedCommandsRule{})
}

type DeprecatedCommandsRule struct{}

func (r *DeprecatedCommandsRule) ID() string { return "GH-010" }
func (r *DeprecatedCommandsRule) Name() string { return "Deprecated Commands" }
func (r *DeprecatedCommandsRule) Category() rules.Category { return rules.CategoryCompliance }
func (r *DeprecatedCommandsRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *DeprecatedCommandsRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *DeprecatedCommandsRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	var findings []rules.Finding
	deprecatedPatterns := []string{"::set-output", "::save-state", "set-output name=", "save-state name="}

	for jobName, job := range pipeline.Jobs {
		for _, step := range job.Steps {
			for _, pattern := range deprecatedPatterns {
				if strings.Contains(step.Run, pattern) {
					findings = append(findings, rules.Finding{
						RuleID:     r.ID(),
						RuleName:   r.Name(),
						Severity:   r.DefaultSeverity(),
						Category:   r.Category(),
						Message:    "Step uses deprecated command pattern: " + pattern,
						Location:   rules.Location{File: pipeline.FilePath, Job: jobName, Line: job.Line},
						Suggestion: "Use environment files (GITHUB_OUTPUT or GITHUB_STATE) instead of echo commands.",
					})
				}
			}
		}
	}
	return findings
}
