package engine

import (
	"testing"

	"github.com/AndreaPallotta/civet/internal/config"
	"github.com/AndreaPallotta/civet/internal/rules"
	_ "github.com/AndreaPallotta/civet/internal/rules/github"
	_ "github.com/AndreaPallotta/civet/internal/rules/gitlab"
	_ "github.com/AndreaPallotta/civet/internal/rules/universal"
)

func TestEngineScoring(t *testing.T) {
	pipeline := &rules.Pipeline{
		Platform: rules.PlatformGitLab,
		FilePath: ".gitlab-ci.yml",
		Stages:   []string{"test"},
		Jobs: map[string]*rules.Job{
			"test": {
				Name:  "test",
				Stage: "test",
				Image: "node:latest", // triggers UNI-002 (Error, -15 on Security)
				Script: []string{
					"npm test",
				},
				Line: 1,
			},
		},
	}

	cfg := config.DefaultConfig()
	report := Analyze(pipeline, cfg)

	if report.OverallScore <= 0 || report.OverallScore > 100 {
		t.Errorf("expected score between 1 and 100, got %d", report.OverallScore)
	}

	if len(report.Findings) == 0 {
		t.Errorf("expected findings to be emitted for mutable image and missing timeout/cache")
	}

	var foundUNI002 bool
	for _, f := range report.Findings {
		if f.RuleID == "UNI-002" {
			foundUNI002 = true
			break
		}
	}
	if !foundUNI002 {
		t.Errorf("expected UNI-002 finding for node:latest")
	}
}

func TestDisabledRules(t *testing.T) {
	pipeline := &rules.Pipeline{
		Platform: rules.PlatformGitLab,
		FilePath: ".gitlab-ci.yml",
		Stages:   []string{"test"},
		Jobs: map[string]*rules.Job{
			"test": {
				Name:   "test",
				Stage:  "test",
				Image:  "node:latest",
				Script: []string{"npm test"},
				Line:   1,
			},
		},
	}

	cfg := &config.Config{
		Rules: config.RulesConfig{
			Disabled: []string{"UNI-002"},
		},
	}

	report := Analyze(pipeline, cfg)
	for _, f := range report.Findings {
		if f.RuleID == "UNI-002" {
			t.Errorf("UNI-002 should have been disabled by config")
		}
	}
}
