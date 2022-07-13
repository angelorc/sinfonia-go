package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagConfig = "config"
)

func addConfigFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagConfig, "./config.yml", "path to config file")
}
