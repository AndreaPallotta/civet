package gitlab

import (
	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&NoTemplateReuseRule{})
}

type NoTemplateReuseRule struct{}

func (r *NoTemplateReuseRule) ID() string { return "GL-004" }
func (r *NoTemplateReuseRule) Name() string { return "No Template Reuse" }
func (r *NoTemplateReuseRule) Category() rules.Category { return rules.CategoryMaintainability }
func (r *NoTemplateReuseRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *NoTemplateReuseRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *NoTemplateReuseRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	if len(pipeline.Jobs) > 5 && len(pipeline.Include) == 0 {
		return []rules.Finding{
			{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Pipeline has many jobs but no 'include' directives.",
				Location:   rules.Location{File: pipeline.FilePath},
				Suggestion: "Use 'include' directives to extract common logic into reusable templates.",
			},
		}
	}

	return nil
}
