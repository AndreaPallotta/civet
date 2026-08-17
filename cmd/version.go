package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.1.1"
	Commit  = "unknown"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of civet",
	Run: func(cmd *cobra.Command, args []string) {
		if Commit != "unknown" && Date != "unknown" {
			fmt.Printf("civet version %s (commit: %s, built: %s)\n", Version, Commit, Date)
		} else {
			fmt.Printf("civet version %s\n", Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
