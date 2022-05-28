package main

import (
	"github.com/angelorc/sinfonia-go/bitsong/cmd/sinfonia-bitsong/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
