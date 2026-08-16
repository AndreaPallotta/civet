package main

import (
	"os"

	"github.com/AndreaPallotta/civet/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
