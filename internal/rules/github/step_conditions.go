package github

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingStepConditionsRule{})
}

type MissingStepConditionsRule struct{}

func (r *MissingStepConditionsRule) ID() string { return "GH-009" }
func (r *MissingStepConditionsRule) Name() string { return "Missing Step Conditions" }
func (r *MissingStepConditionsRule) Category() rules.Category { return rules.CategoryCompliance }
func (r *MissingStepConditionsRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *MissingStepConditionsRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *MissingStepConditionsRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	totalSteps := 0
	stepsWithoutIf := 0

	for _, job := range pipeline.Jobs {
		for _, step := range job.Steps {
			totalSteps++
			if step.If == "" {
				stepsWithoutIf++
			}
		}
	}

	if totalSteps > 0 && float64(stepsWithoutIf)/float64(totalSteps) > 0.70 {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "More than 70% of steps do not have 'if' conditions.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Consider using 'if' conditions to conditionally execute steps and optimize pipeline execution.",
			},
		}
	}

	return nil
}
