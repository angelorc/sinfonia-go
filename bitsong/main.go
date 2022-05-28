package main

import (
	"github.com/angelorc/sinfonia-go/bitsong/cmd/sinfonia-bitsong/cmd"
	"os"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
