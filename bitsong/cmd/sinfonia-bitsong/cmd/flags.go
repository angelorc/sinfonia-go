package cmd

import "github.com/spf13/cobra"

const (
	flagModules     = "modules"
	flagConcurrent  = "concurrent"
	flagStartHeight = "start-height"
	flagEndHeight   = "end-height"
	flagConfig      = "config"
)

func addConfigFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagConfig, "./config.yml", "path to config file")
}
