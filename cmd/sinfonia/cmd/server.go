package cmd

import (
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/server"
	"github.com/spf13/cobra"
	"strconv"
)

func GetServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "module server (rest & graphql)",
	}

	cmd.AddCommand(
		GetServerStartCmd(),
	)

	return cmd
}

func GetServerStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "start server (rest & graphql)",
		Example: "sinfonia server start --mongo-dbname sinfonia-test",
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

			server.Start()

			return nil
		},
	}

	addMongoFlags(cmd)

	return cmd
}
