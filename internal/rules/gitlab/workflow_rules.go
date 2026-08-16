package gitlab

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&NoWorkflowRulesRule{})
}

type NoWorkflowRulesRule struct{}

func (r *NoWorkflowRulesRule) ID() string { return "GL-010" }
func (r *NoWorkflowRulesRule) Name() string { return "No Workflow Rules" }
func (r *NoWorkflowRulesRule) Category() rules.Category { return rules.CategoryCompliance }
func (r *NoWorkflowRulesRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *NoWorkflowRulesRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *NoWorkflowRulesRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	if pipeline.Workflow == nil || len(pipeline.Workflow.Rules) == 0 {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline does not define workflow:rules.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Use workflow:rules to explicitly control when a pipeline is created.",
			},
		}
	}

	return nil
}
