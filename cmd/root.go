package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	format     string
	configFile string
	aiProvider string
)

var rootCmd = &cobra.Command{
	Use:           "civet",
	Short:         "CI/CD pipeline analyzer and assistant",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "terminal", "Output format (terminal/json/markdown)")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", ".civet.yml", "Path to config file")
	rootCmd.PersistentFlags().StringVar(&aiProvider, "ai", "", "AI provider override (claude/openai/gemini/ollama)")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}
