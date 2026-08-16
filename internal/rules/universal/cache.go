package universal

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingCacheRule{})
}

type MissingCacheRule struct{}

func (r *MissingCacheRule) ID() string { return "UNI-001" }
func (r *MissingCacheRule) Name() string { return "Missing Cache Configuration" }
func (r *MissingCacheRule) Category() rules.Category { return rules.CategoryPerformance }
func (r *MissingCacheRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *MissingCacheRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}
func (r *MissingCacheRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	hasCache := false
	for _, job := range pipeline.Jobs {
		if job.Cache != nil {
			hasCache = true
		}
		
		needsCache := false
		for _, s := range job.Script {
			ls := strings.ToLower(s)
			if strings.Contains(ls, "build") || strings.Contains(ls, "install") || strings.Contains(ls, "compile") {
				needsCache = true
			}
		}
		for _, step := range job.Steps {
			ls := strings.ToLower(step.Run)
			if strings.Contains(ls, "build") || strings.Contains(ls, "install") || strings.Contains(ls, "compile") {
				needsCache = true
			}
		}
		if needsCache && job.Cache == nil {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job could benefit from caching (build/install/compile detected)",
				Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
				Suggestion: "Configure a cache for this job to improve performance.",
			})
		}
	}
	
	if !hasCache && len(pipeline.Jobs) > 0 {
		findings = append(findings, rules.Finding{
			RuleID:     r.ID(),
			RuleName:   r.Name(),
			Severity:   r.DefaultSeverity(),
			Category:   r.Category(),
			Message:    "No jobs have a Cache configured in this pipeline",
			Location:   rules.Location{File: pipeline.FilePath},
			Suggestion: "Consider using caching for dependencies and build outputs.",
		})
	}
	return findings
}
