package cmd

import (
	"github.com/angelorc/sinfonia-go/config"
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
			cfgPath, err := config.ParseFlags()
			if err != nil {
				return err
			}

			cfg, err := config.NewConfig(cfgPath)
			if err != nil {
				return err
			}

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             cfg.Mongo.Uri,
				DataBaseName:    cfg.Mongo.DbName,
				RetryWrites:     strconv.FormatBool(cfg.Mongo.Retry),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			server.Start(*cfg)

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}
