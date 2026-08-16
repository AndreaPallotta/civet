package gitlab

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func init() {
	rules.Register(&MissingCoverageReportingRule{})
}

type MissingCoverageReportingRule struct{}

func (r *MissingCoverageReportingRule) ID() string { return "GL-005" }
func (r *MissingCoverageReportingRule) Name() string { return "Missing Coverage Reporting" }
func (r *MissingCoverageReportingRule) Category() rules.Category { return rules.CategoryObservability }
func (r *MissingCoverageReportingRule) DefaultSeverity() rules.Severity { return rules.SeverityInfo }
func (r *MissingCoverageReportingRule) Platforms() []rules.Platform { return []rules.Platform{rules.PlatformGitLab} }

func (r *MissingCoverageReportingRule) Check(pipeline *rules.Pipeline) []rules.Finding {
	if pipeline.Platform != rules.PlatformGitLab {
		return nil
	}

	var findings []rules.Finding
	for name, job := range pipeline.Jobs {
		if strings.Contains(strings.ToLower(name), "test") {
			hasCoverage := false
			if job.Artifacts != nil && job.Artifacts.Reports != nil {
				if _, ok := job.Artifacts.Reports["coverage_report"]; ok { // Note: GitLab uses coverage_report or similar, check logic per requirements
					hasCoverage = true
				}
				// Generic coverage key check
				for k := range job.Artifacts.Reports {
					if strings.Contains(k, "coverage") {
						hasCoverage = true
						break
					}
				}
			}

			if !hasCoverage {
				findings = append(findings, rules.Finding{
					RuleID:     r.ID(),
					RuleName:   r.Name(),
					Severity:   r.DefaultSeverity(),
					Category:   r.Category(),
					Message:    "Test job does not output coverage reports.",
					Location:   rules.Location{File: pipeline.FilePath, Job: name, Line: job.Line},
					Suggestion: "Add an 'artifacts:reports:coverage_report' field to capture test coverage metrics.",
				})
			}
		}
	}
	return findings
}
