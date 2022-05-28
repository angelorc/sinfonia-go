package main

import (
	"github.com/angelorc/sinfonia-go/cmd/sinfonia/cmd"
	"os"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
