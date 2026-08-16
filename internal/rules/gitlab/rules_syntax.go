package gitlab

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&DeprecatedSyntaxRule{})
}

type DeprecatedSyntaxRule struct{}

func (r *DeprecatedSyntaxRule) ID() string { return "GL-001" }
func (r *DeprecatedSyntaxRule) Name() string { return "Deprecated Only/Except Syntax" }
func (r *DeprecatedSyntaxRule) Category() rules.Category { return rules.CategoryCompliance }
func (r *DeprecatedSyntaxRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *DeprecatedSyntaxRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *DeprecatedSyntaxRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		if len(job.Only) > 0 || len(job.Except) > 0 {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job uses deprecated 'only' or 'except' fields.",
				Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
				Suggestion: "Use 'rules' instead of 'only' and 'except' for better control and forward compatibility.",
			})
		}
	}
	return findings
}
