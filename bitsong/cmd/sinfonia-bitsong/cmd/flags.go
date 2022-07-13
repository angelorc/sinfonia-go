package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	flagModules    = "modules"
	flagConcurrent = "concurrent"

	flagStartHeight = "start-height"
	flagEndHeight   = "end-height"

	flagMongoUri    = "mongo-uri"
	flagMongoDBName = "mongo-dbname"
	flagMongoRetry  = "mongo-retry"
)

func addMongoFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagMongoUri, "mongodb://localhost:27017", "the mongo uri connection")
	cmd.Flags().String(flagMongoDBName, "", "the mongo db name to use")
	cmd.Flags().Bool(flagMongoRetry, true, "mongo retrywrites param")
}

func parseMongoFlags(cmd *cobra.Command) (string, string, bool, error) {
	mongoURI, err := cmd.Flags().GetString(flagMongoUri)
	if err != nil || mongoURI == "" {
		return "", "", false, fmt.Errorf("indicate the mongo uri connection\n")
	}

	mongoDBName, err := cmd.Flags().GetString(flagMongoDBName)
	if err != nil || mongoDBName == "" {
		return "", "", false, fmt.Errorf("indicate the mongo db name eg: --mongo-dbname [name]\n")
	}

	mongoRetryWrites, _ := cmd.Flags().GetBool(flagMongoRetry)

	return mongoURI, mongoDBName, mongoRetryWrites, nil
}
