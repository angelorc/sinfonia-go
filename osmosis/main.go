package main

import (
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/osmosis/chain"
	"github.com/angelorc/sinfonia-go/osmosis/indexer"
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

	client, err := chain.NewClient(chain.GetOsmosisConfig())
	if err != nil {
		log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "osmosis", err)
	}

	indexer.NewIndexer(client).Start(2500, 2600, 5)
}
