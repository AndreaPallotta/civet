package universal

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&NoConditionalExecutionRule{})
}

type NoConditionalExecutionRule struct{}

func (r *NoConditionalExecutionRule) ID() string { return "UNI-008" }
func (r *NoConditionalExecutionRule) Name() string { return "No Conditional Execution" }
func (r *NoConditionalExecutionRule) Category() rules.Category { return rules.CategoryPerformance }
func (r *NoConditionalExecutionRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *NoConditionalExecutionRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}

func (r *NoConditionalExecutionRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	hasCondition := false
	
	for _, job := range pipeline.Jobs {
		if len(job.Rules) > 0 || len(job.Only) > 0 || len(job.Except) > 0 {
			hasCondition = true
			break
		}
		if pipeline.Platform == rules.PlatformGitHub {
			for _, step := range job.Steps {
				if step.If != "" {
					hasCondition = true
					break
				}
			}
		}
	}
	
	if !hasCondition && len(pipeline.Jobs) > 0 {
		findings = append(findings, rules.Finding{
			RuleID:     r.ID(),
			RuleName:   r.Name(),
			Severity:   r.DefaultSeverity(),
			Category:   r.Category(),
			Message:    "No jobs specify conditional execution (rules/if)",
			Location:   rules.Location{File: pipeline.FilePath},
			Suggestion: "Use rules or if statements to only run jobs when their relevant files change.",
		})
	}
	return findings
}
