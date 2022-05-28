package sinfonia_bitsong

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sinfonia-bitsong",
		Short: "A bitsong indexer used to collect data for sinfonia",
	}

	// rootCmd.AddCommand( indexer cmds... )

	return rootCmd
}
