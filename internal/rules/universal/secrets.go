package universal

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&HardcodedCredentialsRule{})
}

type HardcodedCredentialsRule struct{}

func (r *HardcodedCredentialsRule) ID() string { return "UNI-005" }
func (r *HardcodedCredentialsRule) Name() string { return "Hardcoded Credentials" }
func (r *HardcodedCredentialsRule) Category() rules.Category { return rules.CategorySecurity }
func (r *HardcodedCredentialsRule) DefaultSeverity() rules.Severity { return rules.SeverityError }
func (r *HardcodedCredentialsRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}

func (r *HardcodedCredentialsRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	patterns := []string{"password=", "token=", "api_key=", "secret=", "aws_secret"}
	
	for _, job := range pipeline.Jobs {
		found := false
		for _, s := range job.Script {
			ls := strings.ToLower(s)
			for _, p := range patterns {
				if strings.Contains(ls, p) {
					found = true
					break
				}
			}
		}
		
		for _, step := range job.Steps {
			ls := strings.ToLower(step.Run)
			for _, p := range patterns {
				if strings.Contains(ls, p) {
					found = true
					break
				}
			}
		}
		
		for k, v := range job.Variables {
			lk := strings.ToLower(k)
			lv := strings.ToLower(v)
			
			// Simple check for base64-like strings (just length for this rule demo)
			if len(v) > 64 && !strings.Contains(v, " ") {
				found = true
			}
			
			for _, p := range patterns {
				if strings.Contains(lk, p) || strings.Contains(lv, p) {
					found = true
					break
				}
			}
		}
		
		if found {
			findings = append(findings, rules.Finding{
				RuleID:     r.ID(),
				RuleName:   r.Name(),
				Severity:   r.DefaultSeverity(),
				Category:   r.Category(),
				Message:    "Job may contain hardcoded credentials",
				Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
				Suggestion: "Use a secure secret management solution or pipeline variables instead of hardcoding secrets.",
			})
		}
	}
	return findings
}
