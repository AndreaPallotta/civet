package github

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingPermissionsRule{})
}

type MissingPermissionsRule struct{}

func (r *MissingPermissionsRule) ID() string { return "GH-002" }
func (r *MissingPermissionsRule) Name() string { return "Missing Permissions Block" }
func (r *MissingPermissionsRule) Category() rules.Category { return rules.CategorySecurity }
func (r *MissingPermissionsRule) DefaultSeverity() rules.Severity { return rules.SeverityError }
func (r *MissingPermissionsRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *MissingPermissionsRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	if len(pipeline.Permissions) == 0 {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline is missing a top-level permissions block.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Explicitly declare 'permissions:' to avoid using the default overly broad GITHUB_TOKEN permissions.",
			},
		}
	}

	return nil
}
