package gitlab

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&ShellExecutorRule{})
}

type ShellExecutorRule struct{}

func (r *ShellExecutorRule) ID() string { return "GL-009" }
func (r *ShellExecutorRule) Name() string { return "Shell Executor Without Sandboxing" }
func (r *ShellExecutorRule) Category() rules.Category { return rules.CategorySecurity }
func (r *ShellExecutorRule) DefaultSeverity() rules.Severity { return rules.SeverityWarning }
func (r *ShellExecutorRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *ShellExecutorRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		hasShellTag := false
		for _, tag := range job.Tags {
			if strings.ToLower(tag) == "shell" {
				hasShellTag = true
				break
			}
		}

		if hasShellTag {
			for _, line := range job.Script {
				if (strings.Contains(line, "curl") || strings.Contains(line, "wget")) &&
					(strings.Contains(line, "| sh") || strings.Contains(line, "| bash")) {
					findings = append(findings, rules.Finding{
						RuleID:     r.ID(),
						RuleName:   r.Name(),
						Severity:   r.DefaultSeverity(),
						Category:   r.Category(),
						Message:    "Job executes remote script directly via shell runner.",
						Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
						Suggestion: "Avoid curling directly to a shell interpreter on un-sandboxed shell runners.",
					})
					break
				}
			}
		}
	}
	return findings
}
