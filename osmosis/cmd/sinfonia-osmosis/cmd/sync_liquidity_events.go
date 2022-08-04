package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math"
	"strconv"
	"strings"
)

func GetSyncLiquidityEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "liquidity-events",
		Short:   "sync liquidity-events (add, exit)",
		Example: "sinfonia-osmosis sync liquidity-events",
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

			if err := syncLiquidityEvents(); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func syncLiquidityEvents() error {
	// get last available height on db
	lastBlock := model.GetLastHeight("osmosis-1")
	// TODO: get first available block
	defaultBlock := 5112889

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.LiquidityEvents = int64(defaultBlock)
	}

	if sync.LiquidityEvents < int64(defaultBlock) {
		sync.LiquidityEvents = int64(defaultBlock)
	}

	txRepo := repository.NewTransactionRepository()
	liquidityRepo := repository.NewLiquidityRepository()
	liquidityRepo.EnsureIndexes()

	limit := 2500
	fromBlock := sync.LiquidityEvents + 1
	toBlock := fromBlock + int64(limit)
	batches := int(math.Ceil(float64(lastBlock-fromBlock) / float64(limit)))

	log.Printf("Scanning blocks from %d to %d, batches %d, first end block %d\n", fromBlock, lastBlock, batches, toBlock)

	for i := 1; i <= batches; i++ {
		if fromBlock > toBlock {
			continue
		}

		events := []bson.M{
			{"events.type": "pool_joined"},
			{"events.type": "pool_exited"},
		}
		txs, err := txRepo.FindEventsByTypes(events, fromBlock, toBlock)
		log.Printf("Scanning blocks from %d to %d, %d txs founds, batch %d/%d\n", fromBlock, toBlock, len(txs), i, batches)

		if err != nil {
			log.Fatalf("Failed to find events. Err: %s", err.Error())
		}

		for _, tx := range txs {
			for _, evt := range tx.Events {
				evtCreate := &modelv2.LiquidityEventCreateReq{
					ChainID:   tx.ChainID,
					Height:    tx.Height,
					TxHash:    tx.Hash,
					Sender:    "",
					TokensIn:  nil,
					TokensOut: nil,
					Time:      tx.Time,
				}

				for _, attr := range evt.Attributes {
					switch attr.Key {
					case "sender":
						evtCreate.Sender = attr.Value
					case "pool_id":
						poolID, _ := strconv.ParseInt(attr.Value, 10, 64)
						evtCreate.PoolID = uint64(poolID)
					case "tokens_in":
						tokensIn, _ := sdk.ParseCoinsNormalized(attr.Value)
						evtCreate.TokensIn = convertCoinsToCoinsModel(tokensIn)
					case "tokens_out":
						tokensOut, _ := sdk.ParseCoinsNormalized(attr.Value)
						evtCreate.TokensOut = convertCoinsToCoinsModel(tokensOut)
					}
				}

				// log.Printf("PoolID: %d, TokensIn: %s, TokensOut: %s", evtCreate.PoolID, evtCreate.TokensIn, evtCreate.TokensOut)
				_, err := liquidityRepo.Create(evtCreate)

				if err != nil {
					if !strings.Contains(err.Error(), "E11000 duplicate key error") {
						log.Fatalf("Failed to write liquidity event to db. Err: %s", err.Error())
					}
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
	sync.LiquidityEvents = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("liquidity events synced to block %d", sync.LiquidityEvents)

	return nil
}
