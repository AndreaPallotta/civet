package gitlab

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&InsecureDinDRule{})
}

type InsecureDinDRule struct{}

func (r *InsecureDinDRule) ID() string { return "GL-008" }
func (r *InsecureDinDRule) Name() string { return "Insecure Docker-in-Docker" }
func (r *InsecureDinDRule) Category() rules.Category { return rules.CategorySecurity }
func (r *InsecureDinDRule) DefaultSeverity() rules.Severity { return rules.SeverityError }
func (r *InsecureDinDRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *InsecureDinDRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		usesDinD := false
		for _, svc := range job.Services {
			if strings.Contains(svc.Name, "docker:dind") {
				usesDinD = true
				break
			}
		}

		if usesDinD {
			tlsVerify := false
			for _, line := range job.Script {
				if strings.Contains(line, "--tls-verify") {
					tlsVerify = true
					break
				}
			}

			tlsCertDir := false
			if val, ok := job.Variables["DOCKER_TLS_CERTDIR"]; ok && val != "" {
				tlsCertDir = true
			}

			if !tlsVerify && !tlsCertDir {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Job uses Docker-in-Docker without TLS enabled.",
					Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
					Suggestion: "Ensure DOCKER_TLS_CERTDIR is set or '--tls-verify' is used when running dind services.",
				})
			}
		}
	}
	return findings
}
