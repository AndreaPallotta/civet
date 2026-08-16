package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AndreaPallotta/civet/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a .civet.yml configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.Generate(".")
		if err != nil {
			return fmt.Errorf("failed to generate config: %w", err)
		}
		fmt.Printf("Successfully created %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
