package universal

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&OverlyBroadTriggerScopeRule{})
}

type OverlyBroadTriggerScopeRule struct{}

func (r *OverlyBroadTriggerScopeRule) ID() string { return "UNI-010" }
func (r *OverlyBroadTriggerScopeRule) Name() string { return "Overly Broad Trigger Scope" }
func (r *OverlyBroadTriggerScopeRule) Category() rules.Category { return rules.CategoryCompliance }
func (r *OverlyBroadTriggerScopeRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *OverlyBroadTriggerScopeRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}

func (r *OverlyBroadTriggerScopeRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	
	if pipeline.Platform == rules.PlatformGitHub {
		if pipeline.Workflow != nil && pipeline.Workflow.Triggers != nil {
			if _, ok := pipeline.Workflow.Triggers["pull_request_target"]; ok {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Workflow is triggered on 'pull_request_target', which can be insecure",
					Location:   rules.Location{File: pipeline.FilePath},
					Suggestion: "Ensure you understand the security implications of 'pull_request_target'.",
				})
			}
		}
	} else if pipeline.Platform == rules.PlatformGitLab {
		if pipeline.Workflow == nil || len(pipeline.Workflow.Rules) == 0 {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Workflow does not define rules, meaning it runs on all triggers by default",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Use 'workflow:rules' to restrict when the pipeline runs.",
			})
		}
	}
	return findings
}
