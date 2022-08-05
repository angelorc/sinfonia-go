package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	"github.com/angelorc/sinfonia-go/osmosis/chain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v9/x/gamm/pool-models/balancer"
	gammtypes "github.com/osmosis-labs/osmosis/v9/x/gamm/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetSyncPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools",
		Short:   "sync pools from latest blocks",
		Example: "sinfonia-osmosis sync pools --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				return err
			}

			cfg, err := config.NewConfig(cfgPath)
			if err != nil {
				return err
			}

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
				return fmt.Errorf("failed to get RPC endpoints on chain %s. err: %v", "osmosis", err)
			}

			if err := syncPools(client); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func GetSyncOldPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "old-pools",
		Short:   "sync old-pools from latest blocks",
		Example: "sinfonia-osmosis sync old-pools ",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				return err
			}

			cfg, err := config.NewConfig(cfgPath)
			if err != nil {
				return err
			}

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
				return fmt.Errorf("failed to get RPC endpoints on chain %s. err: %v", "osmosis", err)
			}

			poolRepo := repository.NewPoolRepository()
			poolRepo.EnsureIndexes()

			// historicalLiqRepo := repository.NewHistoricalLiquidityRepository()

			// defaultBlock := int64(5112879)
			defaultTime := time.Date(2022, 07, 11, 15, 05, 41, 0, time.UTC)

			// import only the first 750 pools, new pools will be imported with the cmd `sync pools`
			for i := 1; i <= 750; i++ {
				// TODO: we need an archive node!!!
				// poolRes, err := client.QueryPoolByIDWithHeight(uint64(i), defaultBlock)
				poolRes, err := client.QueryPoolByID(uint64(i))
				if err != nil {
					return fmt.Errorf("error while fetching poolID, err: %s", err.Error())
				}

				var poolI gammtypes.PoolI
				err = client.Codec.Marshaler.UnpackAny(poolRes.GetPool(), &poolI)
				if err != nil {
					log.Fatalf("error while decoding the new pool")
				}

				pool, ok := poolI.(*balancer.Pool)
				if !ok {
					log.Fatalf("error while decoding the new pool")
				}

				poolAssets := pool.GetAllPoolAssets()

				tracked := false
				inverted := false
				for i, pAsset := range poolAssets {
					if len(poolAssets) == 2 {
						if pAsset.Token.Denom == "uosmo" ||
							pAsset.Token.Denom == "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2" ||
							pAsset.Token.Denom == "ibc/4E5444C35610CC76FC94E7F7886B93121175C28262DDFDDE6F84E82BF2425452" {
							tracked = true
						}

						if pAsset.Token.Denom == "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2" {
							if i == 0 && poolAssets[1].Token.Denom != "uosmo" {
								inverted = true
							}
						}

						if pAsset.Token.Denom == "ibc/4E5444C35610CC76FC94E7F7886B93121175C28262DDFDDE6F84E82BF2425452" {
							if i == 0 && poolAssets[1].Token.Denom != "uosmo" {
								inverted = true
							}
						}
					}
				}

				_, err = poolRepo.Create(&modelv2.PoolCreateReq{
					ChainID:    "osmosis-1",
					Height:     0,
					TxHash:     "",
					PoolID:     uint64(i),
					PoolAssets: convertPoolAssetsToModel(poolAssets),
					SwapFee:    pool.GetSwapFee(sdk.Context{}).MustFloat64(),
					ExitFee:    pool.GetExitFee(sdk.Context{}).MustFloat64(),
					Time:       defaultTime,
					Tracked:    tracked,
					Inverted:   inverted,
				})

				/*if tracked {
					historicalLiqRepo.Create(&modelv2.HistoricalLiquidityCreateReq{
						PoolID: uint64(i),
						Assets: convertPoolAssetsToCoinModel(poolAssets),
						Time:   defaultTime,
					})
				}*/

				if err != nil {
					if !strings.Contains(err.Error(), "E11000 duplicate key error") {
						log.Fatalf("Failed to write pool to db. Err: %s", err.Error())
					}
				}
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func syncPools(client *chain.Client) error {
	// get last available height on db
	lastBlock := model.GetLastHeight("osmosis-1")
	// TODO: get first available block
	defaultBlock := 5112889

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Pools = int64(defaultBlock)
	}

	if sync.Pools < int64(defaultBlock) {
		sync.Pools = int64(defaultBlock)
	}

	txRepo := repository.NewTransactionRepository()
	poolRepo := repository.NewPoolRepository()
	// historicalLiqRepo := repository.NewHistoricalLiquidityRepository()

	limit := 500
	fromBlock := sync.Pools + 1
	toBlock := fromBlock + int64(limit)
	batches := int(math.Ceil(float64(lastBlock-fromBlock) / float64(limit)))

	log.Printf("Scanning blocks from %d to %d, batches %d, first end block %d\n", fromBlock, lastBlock, batches, toBlock)

	for i := 1; i <= batches; i++ {
		if fromBlock > toBlock {
			continue
		}

		log.Printf("Querying blocks from %d to %d", fromBlock, toBlock)
		events := []bson.M{
			{"event.type": "pool_created"},
		}
		txs, err := txRepo.FindEventsByTypes(events, fromBlock, toBlock)
		log.Printf("Scanning blocks from %d to %d, %d txs founds, batch %d/%d\n", fromBlock, toBlock, len(txs), i, batches)

		if err != nil {
			log.Fatalf("Failed to find events. Err: %s", err.Error())
		}

		for _, tx := range txs {
			log.Printf("found %d events", len(tx.Events))

			for _, evt := range tx.Events {
				poolID, err := strconv.ParseUint(evt.Attributes[0].Value, 0, 64)
				if err != nil {
					return fmt.Errorf("error while parsing poolID, err: %s", err.Error())
				}

				poolRes, err := client.QueryPoolByID(poolID)
				if err != nil {
					return fmt.Errorf("error while fetching poolID, err: %s", err.Error())
				}

				var poolI gammtypes.PoolI
				err = client.Codec.Marshaler.UnpackAny(poolRes.GetPool(), &poolI)
				if err != nil {
					log.Fatalf("error while decoding the new pool")
				}

				pool, ok := poolI.(*balancer.Pool)
				if !ok {
					log.Fatalf("error while decoding the new pool")
				}

				poolAssets := pool.GetAllPoolAssets()

				tracked := false
				inverted := false
				for i, pAsset := range poolAssets {
					if len(poolAssets) == 2 {
						if pAsset.Token.Denom == "uosmo" ||
							pAsset.Token.Denom == "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2" ||
							pAsset.Token.Denom == "ibc/4E5444C35610CC76FC94E7F7886B93121175C28262DDFDDE6F84E82BF2425452" {
							tracked = true
						}

						if pAsset.Token.Denom == "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2" {
							if i == 0 && poolAssets[1].Token.Denom != "uosmo" {
								inverted = true
							}
						}

						if pAsset.Token.Denom == "ibc/4E5444C35610CC76FC94E7F7886B93121175C28262DDFDDE6F84E82BF2425452" {
							if i == 0 && poolAssets[1].Token.Denom != "uosmo" {
								inverted = true
							}
						}
					}
				}

				_, err = poolRepo.Create(&modelv2.PoolCreateReq{
					ChainID:    "osmosis-1",
					Height:     tx.Height,
					TxHash:     tx.Hash,
					PoolID:     uint64(i),
					PoolAssets: convertPoolAssetsToModel(pool.GetAllPoolAssets()),
					SwapFee:    pool.GetSwapFee(sdk.Context{}).MustFloat64(),
					ExitFee:    pool.GetExitFee(sdk.Context{}).MustFloat64(),
					Time:       tx.Time,
					Tracked:    tracked,
					Inverted:   inverted,
				})

				/*if tracked {
					historicalLiqRepo.Create(&modelv2.HistoricalLiquidityCreateReq{
						PoolID: uint64(i),
						Assets: convertPoolAssetsToCoinModel(poolAssets),
						Time:   tx.Time,
					})
				}*/

				if err != nil {
					log.Fatalf("Failed to write swap to db. Err: %s", err.Error())
				}

			}
		}

		fromBlock = toBlock + 1
		toBlock = fromBlock + int64(limit)
		if toBlock > lastBlock {
			toBlock = lastBlock
		}
	}

	// update sync with last synced height
	sync.Pools = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("pools synced to block %d", sync.Pools)

	return nil
}

func convertPoolAssetsToModel(pa []balancer.PoolAsset) []modelv2.PoolAsset {
	newPoolAssets := make([]modelv2.PoolAsset, len(pa))

	for i, p := range pa {
		newPoolAssets[i] = modelv2.PoolAsset{
			Token:  modelv2.Coin{Denom: p.Token.Denom, Amount: p.Token.Amount.ToDec().MustFloat64()},
			Weight: p.Weight.String(),
		}
	}

	return newPoolAssets
}

func convertPoolAssetsToCoinModel(pa []balancer.PoolAsset) []modelv2.Coin {
	newAssets := make([]modelv2.Coin, len(pa))

	for i, p := range pa {
		newAssets[i] = modelv2.Coin{Denom: p.Token.Denom, Amount: p.Token.Amount.ToDec().MustFloat64()}
	}

	return newAssets
}
