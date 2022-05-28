package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	flagMongoUri    = "mongo-uri"
	flagMongoDBName = "mongo-dbname"
	flagMongoRetry  = "mongo-retry"
)

func GetAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "module account",
	}

	cmd.AddCommand(GetAccountSyncCmd())

	return cmd
}

func GetAccountSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "sync accounts from latest blocks",
		Example: "sinfonia account sync --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, err := cmd.Flags().GetString(flagMongoUri)
			if err != nil || mongoURI == "" {
				return fmt.Errorf("indicate the mongo uri connection\n")
			}

			mongoDBName, err := cmd.Flags().GetString(flagMongoDBName)
			if err != nil || mongoDBName == "" {
				return fmt.Errorf("indicate the mongo db name eg: --mongo-dbname [name]\n")
			}

			mongoRetryWrites, _ := cmd.Flags().GetBool(flagMongoRetry)

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             mongoURI,
				DataBaseName:    mongoDBName,
				RetryWrites:     strconv.FormatBool(mongoRetryWrites),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			if err := model.SyncAccounts(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagMongoUri, "mongodb://localhost:27017", "the mongo uri connection")
	cmd.Flags().String(flagMongoDBName, "", "the mongo db name to use")
	cmd.Flags().Bool(flagMongoRetry, true, "mongo retrywrites param")

	return cmd
}
