package main

import (
	"github.com/angelorc/sinfonia-go/bitsong/chain"
	"github.com/angelorc/sinfonia-go/bitsong/indexer"
	"log"
)

func main() {
	client, err := chain.NewClient(chain.GetBitsongConfig())
	if err != nil {
		log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
	}

	indexer.NewIndexer(client).Start(1, 10, 5)
}
