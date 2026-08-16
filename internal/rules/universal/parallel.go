package universal

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&SequentialJobsRule{})
}

type SequentialJobsRule struct{}

func (r *SequentialJobsRule) ID() string { return "UNI-007" }
func (r *SequentialJobsRule) Name() string { return "Sequential Jobs Without Dependencies" }
func (r *SequentialJobsRule) Category() rules.Category { return rules.CategoryPerformance }
func (r *SequentialJobsRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *SequentialJobsRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}

func (r *SequentialJobsRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	
	if pipeline.Platform == rules.PlatformGitLab {
		hasNeeds := false
		for _, job := range pipeline.Jobs {
			if len(job.Needs) > 0 {
				hasNeeds = true
				break
			}
		}
		if !hasNeeds && len(pipeline.Stages) > 1 {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline uses sequential stages without declaring any needs dependencies",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Use the 'needs' keyword to enable Directed Acyclic Graph (DAG) for faster execution.",
			})
		}
	} else if pipeline.Platform == rules.PlatformGitHub {
		hasNeeds := false
		for _, job := range pipeline.Jobs {
			if len(job.Needs) > 0 {
				hasNeeds = true
				break
			}
		}
		if hasNeeds {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline relies on 'needs' to serialize jobs. Ensure they actually require sequential execution.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Check if jobs linked by 'needs' can be parallelized.",
			})
		}
	}
	
	return findings
}
