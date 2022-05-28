package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/bitsong/chain"
	"github.com/angelorc/sinfonia-go/bitsong/indexer"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

const (
	flagStartHeight = "start-height"
	flagEndHeight   = "end-height"
	flagConcurrent  = "concurrent"

	flagMongoUri   = "mongo-uri"
	flagMongoDB    = "mongo-db"
	flagMongoRetry = "mongo-retry"
)

func IndexerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "cli to index the bitsong blockchain",
	}

	cmd.AddCommand(indexerStart())

	return cmd
}

func indexerStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "start the indexer from parameters",
		Example: "sinfonia-bitsong indexer start --start-height 1 --end-height 100 --concurrent 2 --mongo-uri \"mongodb://localhost:27017\" --mongo-db sinfonia-test ",
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, err := cmd.Flags().GetString(flagMongoUri)
			if err != nil || mongoURI != "" {
				return fmt.Errorf("indicate the mongo uri connection\n")
			}

			mongoDB, err := cmd.Flags().GetString(flagMongoDB)
			if err != nil && mongoDB != "" {
				return fmt.Errorf("indicate the mongo db\n")
			}

			mongoRetryWrites, _ := cmd.Flags().GetBool(flagMongoRetry)

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             mongoURI,
				DataBaseName:    mongoDB,
				RetryWrites:     strconv.FormatBool(mongoRetryWrites),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			client, err := chain.NewClient(chain.GetBitsongConfig())
			if err != nil {
				log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
			}

			startHeight, err := cmd.Flags().GetInt64(flagStartHeight)
			if err != nil {
				return fmt.Errorf("indicate the start-height\n")
			}

			endHeight, err := cmd.Flags().GetInt64(flagEndHeight)
			if err != nil {
				return fmt.Errorf("indicate the end-height\n")
			}

			concurrent, err := cmd.Flags().GetInt(flagConcurrent)
			if err != nil {
				return fmt.Errorf("indicate the concurrent process\n")
			}

			if startHeight <= 0 {
				return fmt.Errorf("start-height must be > 0\n")
			}

			if startHeight > endHeight {
				return fmt.Errorf("start-height must be < then end-height\n")
			}

			if concurrent > 5 {
				return fmt.Errorf("concurrent is too high\n")
			}

			indexer.NewIndexer(client).Start(startHeight, endHeight, concurrent)

			return nil
		},
	}

	cmd.Flags().Int64(flagStartHeight, 1, "the initial start height")
	cmd.Flags().Int64(flagEndHeight, 100, "the last height to parse")
	cmd.Flags().Int(flagConcurrent, 2, "how many concurrent indexer (do not abuse!)")

	cmd.Flags().String(flagMongoUri, "mongodb://localhost:27017", "the mongo uri connection")
	cmd.Flags().String(flagMongoDB, "", "the mongo db name to use")
	cmd.Flags().Bool(flagMongoRetry, true, "mongo retrywrites param")

	return cmd
}
