package universal

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&DeployWithoutGateRule{})
}

type DeployWithoutGateRule struct{}

func (r *DeployWithoutGateRule) ID() string { return "UNI-009" }
func (r *DeployWithoutGateRule) Name() string { return "Deploy Without Manual Gate" }
func (r *DeployWithoutGateRule) Category() rules.Category { return rules.CategorySecurity }
func (r *DeployWithoutGateRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *DeployWithoutGateRule) Platforms() []rules.Platform {
	return []rules.Platform{rules.PlatformGitLab, rules.PlatformGitHub}
}

func (r *DeployWithoutGateRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	var findings []rules.Finding
	
	for _, job := range pipeline.Jobs {
		lname := strings.ToLower(job.Name)
		if strings.Contains(lname, "deploy") || strings.Contains(lname, "release") || strings.Contains(lname, "publish") {
			hasGate := false
			if pipeline.Platform == rules.PlatformGitLab {
				for _, rule := range job.Rules {
					if rule.When == "manual" {
						hasGate = true
						break
					}
				}
			} else if pipeline.Platform == rules.PlatformGitHub {
				if job.Environment != "" {
					hasGate = true
				}
			}
			
			if !hasGate {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Deploy job lacks a manual approval gate or environment",
					Location:   rules.Location{File: pipeline.FilePath, Job: job.Name, Line: job.Line},
					Suggestion: "Add a manual approval gate or use an environment to protect deployments.",
				})
			}
		}
	}
	return findings
}
