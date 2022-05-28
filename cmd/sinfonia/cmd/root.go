package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sinfonia",
		Short: "A command line to interact with sinfonia backend",
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
		GetAccountCmd(),
	)

	return rootCmd
}
