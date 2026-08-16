package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AndreaPallotta/civet/internal/rules"
)

func TestParseGitLab(t *testing.T) {
	yamlContent := `
stages:
  - test
  - build
  - deploy

variables:
  GLOBAL_VAR: "true"

test_job:
  stage: test
  image: golang:1.22-alpine
  script:
    - go test ./...
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
  cache:
    key: go-mod
    paths:
      - .cache/go-build
  artifacts:
    expire_in: 1 week
    paths:
      - coverage.out
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, ".gitlab-ci.yml")
	if err := os.WriteFile(filePath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	pipeline, err := Parse(filePath)
	if err != nil {
		t.Fatalf("unexpected error parsing gitlab pipeline: %v", err)
	}

	if pipeline.Platform != rules.PlatformGitLab {
		t.Errorf("expected platform GitLab, got %v", pipeline.Platform)
	}

	if len(pipeline.Stages) != 3 {
		t.Errorf("expected 3 stages, got %d", len(pipeline.Stages))
	}

	job, exists := pipeline.Jobs["test_job"]
	if !exists {
		t.Fatalf("expected test_job to exist")
	}

	if job.Image != "golang:1.22-alpine" {
		t.Errorf("expected image golang:1.22-alpine, got %s", job.Image)
	}

	if job.Cache == nil || job.Cache.Key != "go-mod" {
		t.Errorf("expected cache key 'go-mod'")
	}
}

func TestParseGitHub(t *testing.T) {
	yamlContent := `
name: CI
on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

concurrency:
  group: ${{ github.ref }}
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@v4
      - name: Run Build
        run: go build ./...
`
	tmpDir := t.TempDir()
	workflowDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	filePath := filepath.Join(workflowDir, "ci.yml")
	if err := os.WriteFile(filePath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	pipeline, err := Parse(filePath)
	if err != nil {
		t.Fatalf("unexpected error parsing github workflow: %v", err)
	}

	if pipeline.Platform != rules.PlatformGitHub {
		t.Errorf("expected platform GitHub, got %v", pipeline.Platform)
	}

	if pipeline.Concurrency == nil || pipeline.Concurrency.Group == "" {
		t.Errorf("expected concurrency group to be parsed")
	}

	job, exists := pipeline.Jobs["build"]
	if !exists {
		t.Fatalf("expected job 'build' to exist")
	}

	if job.RunsOn != "ubuntu-latest" {
		t.Errorf("expected runs-on ubuntu-latest, got %s", job.RunsOn)
	}

	if len(job.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(job.Steps))
	}
}
