package cmd

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/indexer"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	"github.com/angelorc/sinfonia-go/osmosis/chain"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

func IndexerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "cli to index the osmosis blockchain",
	}

	cmd.AddCommand(GetIndexerParserCmd())

	return cmd
}

func GetIndexerParserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse [from-block] [to-block]",
		Short:   "parse the osmosis blockchain from block to block",
		Example: "sinfonia-osmosis indexer parse 1 100 --concurrent 2 --mongo-dbname sinfonia-test",
		Args:    cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString(flagConfig)
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

			client, err := chain.NewClient(&cfg.Osmosis)
			if err != nil {
				log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "osmosis", err)
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

			syncAll := false

			if startHeight <= 0 {
				startHeight = blockRepo.Latest().Height + 1
				syncAll = true
			}

			if endHeight <= startHeight {
				endHeight = client.LatestBlockHeight(context.Background())
				syncAll = true
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

			if syncAll {
				if err := syncPools(client); err != nil {
					return err
				}

				if err := syncLiquidityEvents(); err != nil {
					return err
				}

				if err := syncSwaps(); err != nil {
					return err
				}
			}

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
