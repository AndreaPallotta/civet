package gitlab

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingResourceGroupRule{})
}

type MissingResourceGroupRule struct{}

func (r *MissingResourceGroupRule) ID() string { return "GL-006" }
func (r *MissingResourceGroupRule) Name() string { return "Missing Resource Group on Deploy" }
func (r *MissingResourceGroupRule) Category() rules.Category { return rules.CategoryReliability }
func (r *MissingResourceGroupRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingResourceGroupRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *MissingResourceGroupRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		lowerName := strings.ToLower(name)
		if strings.Contains(lowerName, "deploy") || strings.Contains(lowerName, "release") || strings.Contains(lowerName, "publish") {
			if job.ResourceGroup == "" {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Deploy job is missing a 'resource_group'.",
					Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
					Suggestion: "Add 'resource_group' to prevent concurrent deployments to the same environment.",
				})
			}
		}
	}
	return findings
}
