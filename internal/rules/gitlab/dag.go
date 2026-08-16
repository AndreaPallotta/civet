package gitlab

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingDAGOptimizationRule{})
}

type MissingDAGOptimizationRule struct{}

func (r *MissingDAGOptimizationRule) ID() string { return "GL-002" }
func (r *MissingDAGOptimizationRule) Name() string { return "Missing DAG Optimization" }
func (r *MissingDAGOptimizationRule) Category() rules.Category { return rules.CategoryPerformance }
func (r *MissingDAGOptimizationRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingDAGOptimizationRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *MissingDAGOptimizationRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	if len(pipeline.Stages) <= 2 {
		return nil
	}

	for _, job := range pipeline.Jobs {
		if len(job.Needs) > 0 {
			return nil
		}
	}

	return []rules.Finding{
		{
			RuleID:     r.ID(),
			RuleName:   r.Name(),
			Severity:   r.DefaultSeverity(),
			Category:   r.Category(),
			Message:    "Pipeline has more than 2 stages but no job uses the 'needs' keyword.",
			Location:   rules.Location{File: pipeline.FilePath},
			Suggestion: "Use Directed Acyclic Graphs (DAG) via the 'needs' keyword to parallelize jobs and speed up the pipeline.",
		},
	}
}
