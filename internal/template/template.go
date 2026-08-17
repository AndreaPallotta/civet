package template

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates
var FS embed.FS

type Options struct {
	WithDocker bool
}

// Render loads and renders a CI/CD pipeline template.
func Render(platform, language string, opts Options) (string, error) {
	path := fmt.Sprintf("templates/%s/%s.yml", platform, language)
	
	content, err := FS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template not found for %s %s: %w", platform, language, err)
	}

	tmpl, err := template.New("pipeline").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, opts); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return strings.TrimSpace(buf.String()) + "\n", nil
}
