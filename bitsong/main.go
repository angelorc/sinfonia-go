package main

import (
	"github.com/angelorc/sinfonia-go/bitsong/chain"
	"github.com/angelorc/sinfonia-go/bitsong/indexer"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"log"
)

func main() {
	/**
	 * Connect to db
	 */
	defaultDB := db.Database{
		DataBaseRefName: "default",
		URL:             config.GetSecret("MONGO_URI"),
		DataBaseName:    config.GetSecret("MONGO_DATABASE"),
		RetryWrites:     config.GetSecret("MONGO_RETRYWRITES"),
	}
	defaultDB.Init()
	defer defaultDB.Disconnect()

	client, err := chain.NewClient(chain.GetBitsongConfig())
	if err != nil {
		log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
	}

	indexer.NewIndexer(client).Parse(1, 1000, 10)
}
