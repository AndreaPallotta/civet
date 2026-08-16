package github

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&ActionPinningRule{})
}

type ActionPinningRule struct{}

func (r *ActionPinningRule) ID() string { return "GH-001" }
func (r *ActionPinningRule) Name() string { return "Actions Not Pinned to SHA" }
func (r *ActionPinningRule) Category() rules.Category { return rules.CategorySecurity }
func (r *ActionPinningRule) DefaultSeverity() rules.Severity { return rules.SeverityError }
func (r *ActionPinningRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *ActionPinningRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	var findings []rules.Finding
	for jobName, job := range pipeline.Jobs {
		for _, step := range job.Steps {
			if step.Uses == "" {
				continue
			}

			// Exception for 'actions/' org
			if strings.HasPrefix(step.Uses, "actions/") {
				continue
			}

			parts := strings.SplitN(step.Uses, "@", 2)
			if len(parts) != 2 {
				continue
			}

			version := parts[1]
			// Check if version is 40 hex characters
			if len(version) != 40 || !isHex(version) {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Action is not pinned to a full SHA.",
					Location:   rules.Location{File: pipeline.FilePath, Job: jobName, Line: job.Line},
					Suggestion: "Pin actions to a full 40-character SHA instead of a tag or branch for better security.",
				})
			}
		}
	}
	return findings
}

func isHex(s string) bool {
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}
