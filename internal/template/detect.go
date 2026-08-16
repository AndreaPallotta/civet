package template

import (
	"os"
	"path/filepath"
)

type DetectionResult struct {
	Language  string
	Framework string
	HasDocker bool
	Platform  string
}

// Detect inspects a directory to guess the language, platform, and docker usage.
func Detect(dir string) (DetectionResult, error) {
	var res DetectionResult

	if exists(filepath.Join(dir, "go.mod")) {
		res.Language = "go"
	} else if exists(filepath.Join(dir, "package.json")) {
		res.Language = "node"
	} else if exists(filepath.Join(dir, "Cargo.toml")) {
		res.Language = "rust"
	} else if exists(filepath.Join(dir, "pyproject.toml")) || exists(filepath.Join(dir, "requirements.txt")) {
		res.Language = "python"
	}

	if exists(filepath.Join(dir, "Dockerfile")) {
		res.HasDocker = true
	}

	if exists(filepath.Join(dir, ".gitlab-ci.yml")) {
		res.Platform = "gitlab"
	} else {
		// check if .github directory exists
		if stat, err := os.Stat(filepath.Join(dir, ".github")); err == nil && stat.IsDir() {
			res.Platform = "github"
		}
	}

	return res, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
