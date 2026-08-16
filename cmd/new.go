package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/AndreaPallotta/civet/internal/template"
)

var (
	scaffoldPlatform string
	scaffoldLang     string
	scaffoldWith     string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Scaffold a new CI/CD pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := template.Detect(".")
		if err != nil {
			fmt.Printf("Warning: failed to auto-detect environment: %v\n", err)
		}

		plat := scaffoldPlatform
		if plat == "" {
			plat = res.Platform
		}
		if plat == "" {
			plat = "gitlab" // fallback default
		}

		lang := scaffoldLang
		if lang == "" {
			lang = res.Language
		}
		if lang == "" {
			lang = "go" // fallback default
		}

		opts := template.Options{
			WithDocker: res.HasDocker || scaffoldWith == "docker",
		}

		content, err := template.Render(plat, lang, opts)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		var dest string
		if plat == "github" {
			dest = ".github/workflows/ci.yml"
			if err := os.MkdirAll(".github/workflows", 0755); err != nil {
				return err
			}
		} else {
			dest = ".gitlab-ci.yml"
		}

		if err := os.WriteFile(dest, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write pipeline file: %w", err)
		}

		fmt.Printf("Successfully generated pipeline for %s (%s) at %s\n", lang, plat, dest)
		fmt.Printf("Suggestion: run `civet audit %s` to check the generated file.\n", dest)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVar(&scaffoldPlatform, "platform", "", "Target platform (gitlab/github)")
	newCmd.Flags().StringVar(&scaffoldLang, "lang", "", "Target language (go/node/python/rust)")
	newCmd.Flags().StringVar(&scaffoldWith, "with", "", "Additional integrations (e.g. docker)")
	rootCmd.AddCommand(newCmd)
}
