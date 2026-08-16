package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
	Commit  = "unknown"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of civet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("civet version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
