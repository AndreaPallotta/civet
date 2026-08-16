package github

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&UnsafePullRequestTargetRule{})
}

type UnsafePullRequestTargetRule struct{}

func (r *UnsafePullRequestTargetRule) ID() string { return "GH-006" }
func (r *UnsafePullRequestTargetRule) Name() string { return "Unsafe pull_request_target" }
func (r *UnsafePullRequestTargetRule) Category() rules.Category { return rules.CategorySecurity }
func (r *UnsafePullRequestTargetRule) DefaultSeverity() rules.Severity { return rules.SeverityError }
func (r *UnsafePullRequestTargetRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *UnsafePullRequestTargetRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	if pipeline.Workflow != nil && pipeline.Workflow.Triggers != nil {
		if _, ok := pipeline.Workflow.Triggers["pull_request_target"]; ok {
			return []rules.Finding{
				{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Workflow uses 'pull_request_target' trigger.",
					Location:   rules.Location{File: pipeline.FilePath},
					Suggestion: "Ensure you are not executing untrusted code from PRs, as this trigger has write access to the base repository.",
				},
			}
		}
	}

	return nil
}
