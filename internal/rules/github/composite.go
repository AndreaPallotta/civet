package github

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&NoCompositeActionReuseRule{})
}

type NoCompositeActionReuseRule struct{}

func (r *NoCompositeActionReuseRule) ID() string { return "GH-004" }
func (r *NoCompositeActionReuseRule) Name() string { return "No Composite Action Reuse" }
func (r *NoCompositeActionReuseRule) Category() rules.Category { return rules.CategoryMaintainability }
func (r *NoCompositeActionReuseRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *NoCompositeActionReuseRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitHub} }

func (r *NoCompositeActionReuseRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitHub {
		return nil
	}

	if len(pipeline.Jobs) <= 5 {
		return nil
	}

	stepCounts := make(map[string]map[string]bool)

	for jobName, job := range pipeline.Jobs {
		for _, step := range job.Steps {
			key := ""
			if step.Uses != "" {
				key = "uses:" + step.Uses
			} else if step.Run != "" {
				key = "run:" + step.Run
			}

			if key != "" {
				if stepCounts[key] == nil {
					stepCounts[key] = make(map[string]bool)
				}
				stepCounts[key][jobName] = true
			}
		}
	}

	hasDuplicates := false
	for _, jobs := range stepCounts {
		if len(jobs) > 1 {
			hasDuplicates = true
			break
		}
	}

	if hasDuplicates {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline has >5 jobs and duplicated steps across different jobs.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Refactor duplicated steps into local composite actions to improve maintainability.",
			},
		}
	}

	return nil
}
