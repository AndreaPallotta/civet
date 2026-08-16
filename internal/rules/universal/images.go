package universal

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MutableImageTagsRule{})
}

type MutableImageTagsRule struct{}

func (r *MutableImageTagsRule) ID() string { return "UNI-002" }
func (r *MutableImageTagsRule) Name() string { return "Mutable Image Tags" }
func (r *MutableImageTagsRule) Category() rules.Category { return rules.CategorySecurity }
func (r *MutableImageTagsRule) DefaultSeverity() rules.Severity { return rules.SeverityError }
func (r *MutableImageTagsRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}
func (r *MutableImageTagsRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	badTags := []string{":latest", ":main", ":master", ":dev"}
	
	for _, job := range pipeline.Jobs {
		if job.Image != "" {
			bad := false
			if !strings.Contains(job.Image, ":") {
				bad = true // No tag provided
			} else {
				for _, t := range badTags {
					if strings.HasSuffix(strings.ToLower(job.Image), t) {
						bad = true
						break
					}
				}
			}
			
			if bad {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Job uses a mutable image tag or no tag",
					Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
					Suggestion: "Use a specific immutable tag or digest for the image to ensure reproducibility.",
				})
			}
		}
	}
	return findings
}
