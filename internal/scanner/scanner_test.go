package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	glFile := filepath.Join(tmpDir, ".gitlab-ci.yml")
	if err := os.WriteFile(glFile, []byte("stages:\n  - test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ghDir := filepath.Join(tmpDir, ".github", "workflows")
	if err := os.MkdirAll(ghDir, 0755); err != nil {
		t.Fatal(err)
	}
	ghFile := filepath.Join(ghDir, "ci.yml")
	if err := os.WriteFile(ghFile, []byte("name: CI\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Ignored folder
	nodeModules := filepath.Join(tmpDir, "node_modules", ".github", "workflows")
	if err := os.MkdirAll(nodeModules, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nodeModules, "ignored.yml"), []byte("name: Ignored\n"), 0644); err != nil {
		t.Fatal(err)
	}

	found, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("expected 2 files discovered, got %d: %v", len(found), found)
	}
}
