package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sinfonia-osmosis",
		Short: "An osmosis indexer used to collect data for sinfonia",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: false,
			HiddenDefaultCmd:    true,
		},
	}

	rootCmd.SetHelpCommand(
		&cobra.Command{
			Hidden: true,
		},
	)

	rootCmd.AddCommand(
		IndexerCmd(),
		GetSyncCmd(),
	)

	return rootCmd
}
