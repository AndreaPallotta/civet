package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
)

// Parse reads the given file, detects the platform based on the path,
// and delegates to the appropriate platform parser.
func Parse(filePath string) (*rules.Pipeline, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	content := string(data)

	// Detect platform based on file path
	base := filepath.Base(filePath)
	dir := filepath.ToSlash(filepath.Dir(filePath))

	if base == ".gitlab-ci.yml" || strings.HasSuffix(base, ".gitlab-ci.yml") {
		return ParseGitLab(content, filePath)
	}

	if strings.Contains(dir, ".github/workflows") && (strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".yaml")) {
		return ParseGitHub(content, filePath)
	}

	return nil, fmt.Errorf("could not detect platform for file: %s", filePath)
}
