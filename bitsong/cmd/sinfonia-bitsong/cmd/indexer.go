package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/bitsong/chain"
	"github.com/angelorc/sinfonia-go/indexer"
	"github.com/angelorc/sinfonia-go/mongo/db"
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
		Use:     "parse [from-block] [to-block]",
		Short:   "parse the bitsong blockchain from block to block",
		Example: "sinfonia-bitsong indexer parse 1 100 --concurrent 2 --mongo-dbname sinfonia-test",
		Args:    cobra.RangeArgs(0, 2),
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

			client, err := chain.NewClient(chain.GetBitsongConfig())
			if err != nil {
				log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
			}

			fromBlock := int64(1)
			toBlock := fromBlock

			if len(args) > 0 {
				fromBlock, err = strconv.ParseInt(args[0], 0, 64)
			}

			if len(args) > 1 {
				toBlock, err = strconv.ParseInt(args[1], 0, 64)
			}

			if fromBlock <= 0 {
				fromBlock = 1
			}

			if toBlock <= fromBlock {
				toBlock = fromBlock
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
				Parse(fromBlock, toBlock)

			return nil
		},
	}

	cmd.Flags().String(flagModules, "*", "modules to parse eg: * for all or \"blocks,transactions,messages,block-results\" ")
	cmd.Flags().Int(flagConcurrent, 2, "how many concurrent indexer (do not abuse!)")

	cmd.Flags().String(flagMongoUri, "mongodb://localhost:27017", "the mongo uri connection")
	cmd.Flags().String(flagMongoDBName, "", "the mongo db name to use")
	cmd.Flags().Bool(flagMongoRetry, true, "mongo retrywrites param")

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
