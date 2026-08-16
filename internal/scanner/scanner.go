package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// Scan recursively finds CI/CD pipeline files in the given directory.
func Scan(dir string) ([]string, error) {
	var results []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || name == "vendor" || name == ".git" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		base := filepath.Base(path)
		if base == ".gitlab-ci.yml" {
			results = append(results, path)
		} else if strings.HasPrefix(filepath.ToSlash(path), ".github/workflows/") && strings.HasSuffix(base, ".yml") {
			results = append(results, path)
		} else if strings.Contains(filepath.ToSlash(path), "/.github/workflows/") && strings.HasSuffix(base, ".yml") {
			results = append(results, path)
		}

		return nil
	})

	return results, err
}
