package cmd

import (
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/spf13/cobra"
	"strconv"
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
			mongoURI, mongoDBName, mongoRetryWrites, err := parseMongoFlags(cmd)
			if err != nil {
				return err
			}

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

	addMongoFlags(cmd)

	return cmd
}