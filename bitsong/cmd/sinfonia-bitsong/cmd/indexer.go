package cmd

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/bitsong/chain"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/indexer"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

func IndexerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "cli to index the bitsong blockchain",
	}

	cmd.AddCommand(GetIndexerParserCmd())

	return cmd
}

func GetIndexerParserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse",
		Short:   "parse the bitsong blockchain from block to block",
		Example: "sinfonia-bitsong indexer parse --start-height=1 --end-height=100 --concurrent 2 --mongo-dbname sinfonia-test",
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

			client, err := chain.NewClient(&cfg.Bitsong)
			if err != nil {
				log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
			}

			blockRepo := repository.NewBlockRepository()

			startHeight, err := cmd.Flags().GetInt64(flagStartHeight)
			if err != nil {
				return err
			}

			endHeight, err := cmd.Flags().GetInt64(flagEndHeight)
			if err != nil {
				return err
			}

			if startHeight <= 0 {
				startHeight = blockRepo.Latest().Height + 1
			}

			if endHeight <= startHeight {
				endHeight = client.LatestBlockHeight(context.Background())
			}

			concurrent, err := cmd.Flags().GetInt(flagConcurrent)
			if err != nil {
				return fmt.Errorf("indicate the concurrent process\n")
			}

			if concurrent > 5 {
				return fmt.Errorf("concurrent is too high\n")
			}

			modulesStr, err := cmd.Flags().GetString(flagModules)
			if err != nil {
				return fmt.Errorf("indicate modules to parse")
			}

			indexer.
				NewIndexer(client, parseModules(modulesStr), concurrent).
				Parse(startHeight, endHeight)

			return nil
		},
	}

	cmd.Flags().Int64(flagStartHeight, 0, "parse from height, default is 0 mean that will be used the latest block stored in db")
	cmd.Flags().Int64(flagEndHeight, 0, "parse to height, default is 0 mean that will be used the current block on chain")

	cmd.Flags().String(flagModules, "*", "modules to parse eg: * for all or \"blocks,transactions,messages,block-results\" ")
	cmd.Flags().Int(flagConcurrent, 2, "how many concurrent indexer (do not abuse!)")

	addConfigFlag(cmd)

	return cmd
}

func parseModules(flag string) *indexer.IndexModules {
	modulesStr := strings.Split(flag, ",")
	modules := &indexer.IndexModules{}

	if modulesStr[0] == "*" {
		modules.Blocks = true
		modules.Transactions = true
		modules.Messages = true
		modules.BlockResults = true

		return modules
	}

	for _, m := range modulesStr {
		switch m {
		case "blocks":
			modules.Blocks = true
		case "transactions":
			modules.Transactions = true
		case "messages":
			modules.Messages = true
		case "block-results":
			modules.BlockResults = true
		}
	}

	return modules
}
