package gitlab

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingInterruptibleRule{})
}

type MissingInterruptibleRule struct{}

func (r *MissingInterruptibleRule) ID() string { return "GL-003" }
func (r *MissingInterruptibleRule) Name() string { return "Missing Interruptible Flag" }
func (r *MissingInterruptibleRule) Category() rules.Category { return rules.CategoryReliability }
func (r *MissingInterruptibleRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *MissingInterruptibleRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *MissingInterruptibleRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		lowerName := strings.ToLower(name)
		isDeploy := strings.Contains(lowerName, "deploy") || strings.Contains(lowerName, "release") || strings.Contains(lowerName, "publish")
		
		if !isDeploy && job.Interruptible == nil {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Non-deploy job is missing the 'interruptible' flag.",
				Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
				Suggestion: "Set 'interruptible: true' to allow safe cancellation of redundant pipeline runs.",
			})
		}
	}
	return findings
}
