package cmd

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	"github.com/angelorc/sinfonia-go/osmosis/chain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v9/x/incentives/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "module sync",
	}

	cmd.AddCommand(
		GetSyncOldPoolCmd(),
		GetSyncPoolCmd(),
		GetSyncSwapCmd(),
		GetSyncIncentivesCmd(),
		GetSyncPricesCmd(),
		GetSyncHistoricalPricesCmd(),
		GetSyncLiquidityEventsCmd(),
	)

	return cmd
}

func GetSyncSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swaps",
		Short:   "sync swaps from latest blocks",
		Example: "sinfonia-osmosis sync swaps --mongo-dbname sinfonia-test",
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

			if err := syncSwaps(); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func chunkSlice(slice []modelv2.Attribute, chunkSize int) [][]modelv2.Attribute {
	var chunks [][]modelv2.Attribute
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func convertCoinToCoinModel(coin sdk.Coin) modelv2.Coin {
	return modelv2.Coin{
		Amount: coin.Amount.ToDec().MustFloat64(),
		Denom:  coin.Denom,
	}
}

func convertCoinsToCoinsModel(coins []sdk.Coin) []modelv2.Coin {
	output := make([]modelv2.Coin, 0)

	for _, coin := range coins {
		outputCoin := convertCoinToCoinModel(coin)
		output = append(output, outputCoin)
	}

	return output
}

func syncSwaps() error {
	// get last available height on db
	lastBlock := model.GetLastHeight("osmosis-1")
	// TODO: get first available block
	defaultBlock := 5112889

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Swaps = int64(defaultBlock)
	}

	if sync.Swaps < int64(defaultBlock) {
		sync.Swaps = int64(defaultBlock)
	}

	txRepo := repository.NewTransactionRepository()
	swapRepo := repository.NewSwapRepository()
	poolRepo := repository.NewPoolRepository()
	hpr := repository.NewHistoricalPriceRepository()

	limit := 2000
	fromBlock := sync.Swaps + 1
	toBlock := fromBlock + int64(limit)
	batches := int(math.Ceil(float64(lastBlock-fromBlock) / float64(limit)))

	log.Printf("Scanning blocks from %d to %d, batches %d, first end block %d\n", fromBlock, lastBlock, batches, toBlock)

	for i := 1; i <= batches; i++ {
		if fromBlock > toBlock {
			continue
		}

		events := []bson.M{
			{"events.type": "token_swapped"},
		}
		txs, err := txRepo.FindEventsByTypes(events, fromBlock, toBlock)
		log.Printf("Scanning blocks from %d to %d, %d txs founds, batch %d/%d\n", fromBlock, toBlock, len(txs), i, batches)

		if err != nil {
			log.Fatalf("Failed to find events. Err: %s", err.Error())
		}

		for _, tx := range txs {
			for _, evt := range tx.Events {
				groupedAttrs := chunkSlice(evt.Attributes, 5)

				for _, attrs := range groupedAttrs {
					swapCreate := &modelv2.SwapCreateReq{
						ID:       primitive.NewObjectID(),
						ChainID:  tx.ChainID,
						Height:   tx.Height,
						TxHash:   tx.Hash,
						Fee:      0,
						UsdValue: 0,
						Time:     tx.Time,
					}

					for _, attr := range attrs {
						switch attr.Key {
						case "sender":
							swapCreate.Account = attr.Value
						case "pool_id":
							poolID, _ := strconv.ParseInt(attr.Value, 10, 64)
							swapCreate.PoolId = poolID
						case "tokens_in":
							tokenIn, _ := sdk.ParseCoinNormalized(attr.Value)
							swapCreate.TokenIn = convertCoinToCoinModel(tokenIn)
						case "tokens_out":
							tokenOut, _ := sdk.ParseCoinNormalized(attr.Value)
							swapCreate.TokenOut = convertCoinToCoinModel(tokenOut)
						}
					}

					pool := poolRepo.FindByPoolID(uint64(swapCreate.PoolId))
					if pool.SwapFee > 0 {
						swapCreate.Fee = calcFee(swapCreate.TokenIn.String(), pool.SwapFee)
					}

					if pool.Tracked {
						if pool.GetBaseAsset().Denom == swapCreate.TokenIn.Denom {
							swapCreate.Type = 0 // buy
						} else {
							swapCreate.Type = 1 // sell
						}

						/*if swapCreate.Type == 0 {
							swapCreate.PriceBase = swapCreate.TokenIn.Amount / swapCreate.TokenOut.Amount
							swapCreate.PriceQuote = swapCreate.TokenOut.Amount / swapCreate.TokenIn.Amount
						} else {
							swapCreate.PriceBase = swapCreate.TokenOut.Amount / swapCreate.TokenIn.Amount
							swapCreate.PriceQuote = swapCreate.TokenIn.Amount / swapCreate.TokenOut.Amount
						}*/

						// add usd value
						price := float64(0)
						prices := hpr.FindByAsset(pool.GetQuoteAsset().Denom, tx.Time)
						if len(prices) > 0 {
							price = prices[0].Price
						}

						if swapCreate.Type == 0 {
							swapCreate.UsdValue = (swapCreate.TokenOut.Amount * 0.000001) * price
						} else {
							swapCreate.UsdValue = (swapCreate.TokenIn.Amount * 0.000001) * price
						}

						// save swap
						//swapCreateBatch = append(swapCreateBatch, swapCreate)
						_, err := swapRepo.Create(swapCreate)

						if err != nil {
							if !strings.Contains(err.Error(), "E11000 duplicate key error") {
								log.Fatalf("Failed to write swap to db. Err: %s", err.Error())
							}
						}

						/*price := swapCreate.UsdValue / (swapCreate.TokenIn.Amount * 0.000001)
						historyCreateBatch = append(historyCreateBatch, &modelv2.HistoricalPriceCreateReq{
							Asset: swapCreate.TokenIn.Denom,
							Price: price,
							Time:  tx.Time,
						})

						price = swapCreate.UsdValue / (swapCreate.TokenOut.Amount * 0.000001)
						historyCreateBatch = append(historyCreateBatch, &modelv2.HistoricalPriceCreateReq{
							Asset: swapCreate.TokenOut.Denom,
							Price: price,
							Time:  tx.Time,
						})*/
					}
				}
			}
		}

		// update sync with last synced height
		sync.Swaps = toBlock
		if err := sync.Save(); err != nil {
			return err
		}

		fromBlock = toBlock + 1
		toBlock = fromBlock + int64(limit)
		if toBlock > lastBlock {
			toBlock = lastBlock
		}
	}

	fmt.Printf("swaps synced to block %d", sync.Swaps)

	return nil
}

func calcFee(tokenInStr string, swapFee float64) float64 {
	tokenIn, _ := sdk.ParseCoinNormalized(tokenInStr)
	swapFeeDec := sdk.MustNewDecFromStr(fmt.Sprintf("%f", swapFee))
	tokenInAfterFee := tokenIn.Amount.ToDec().Mul(sdk.OneDec().Sub(swapFeeDec)).TruncateInt()

	return tokenIn.Amount.Sub(tokenInAfterFee).ToDec().MustFloat64()
}

func calcVolumeUSD(tokensIn, tokensOut string, ts time.Time) float64 {
	ibcDenom := "ibc/8B066EED78CCC6A90E963C81EB4B527C28FE538BE396B8756F4C4BFC53C74221"
	amtBTSG := int64(0)

	if strings.HasSuffix(tokensIn, ibcDenom) {
		amt := strings.Replace(tokensIn, ibcDenom, "", -1)
		amtBTSG, _ = strconv.ParseInt(amt, 0, 64)
	}

	if strings.HasSuffix(tokensOut, ibcDenom) {
		amt := strings.Replace(tokensOut, ibcDenom, "", -1)
		amtBTSG, _ = strconv.ParseInt(amt, 0, 64)
	}

	prices := make(map[string]float64)
	prices["09-05-2022"] = 0.04840136619600178
	prices["10-05-2022"] = 0.039866432548801504
	prices["11-05-2022"] = 0.03766485405081663
	prices["12-05-2022"] = 0.033936829238417024
	prices["13-05-2022"] = 0.02447988666488955
	prices["14-05-2022"] = 0.022366606595902307
	prices["15-05-2022"] = 0.02193695925498291
	prices["16-05-2022"] = 0.025348193585375472
	prices["17-05-2022"] = 0.02257440395770441
	prices["18-05-2022"] = 0.02299961941622652
	prices["19-05-2022"] = 0.020301519267571042

	volume := (float64(amtBTSG) * 0.000001) * prices[ts.Format("02-01-2006")]

	//return fmt.Sprintf("%f", volume)
	return volume
}

func GetSyncIncentivesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "incentives",
		Short:   "sync incentives from latest blocks",
		Example: "sinfonia-osmosis sync incentives",
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

			if err := syncIncentives(client); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func syncIncentives(client *chain.Client) error {
	// get last available height on db
	lastBlock := model.GetLastHeight("osmosis-1")
	// TODO: get first available block
	defaultBlock := 5112889

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Incentives = int64(defaultBlock)
	}

	if sync.Incentives < int64(defaultBlock) {
		sync.Incentives = int64(defaultBlock)
	}

	fromBlock := sync.Incentives + 1

	for height := fromBlock; height < lastBlock; height++ {
		incentiveRepo := repository.NewIncentiveRepository()

		log.Printf("querying block results, height %d", height)

		ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)

		blockResults, err := client.QueryBlockResults(ctx, &height)
		if err != nil {
			return fmt.Errorf("error while fetching blockresults, err: %s", err.Error())
		}

		log.Printf("iterating block results, height %d", height)

		for _, evt := range blockResults.BeginBlockEvents {
			switch evt.Type {
			case types.TypeEvtDistribution:
				incentive := modelv2.IncentiveCreateReq{
					ChainID: "osmosis-1",
					Height:  height,
					Time:    time.Now(), // add block time
				}

				for _, attr := range evt.Attributes {
					switch string(attr.Key) {
					case types.AttributeReceiver:
						incentive.Receiver = string(attr.Value)
					case types.AttributeAmount:
						assets, err := sdk.ParseCoinsNormalized(string(attr.Value))
						if err != nil {
							log.Fatalf("error while converting coins")
						}
						incentive.Assets = convertCoinsToCoinsModel(assets)
					}
				}

				_, err := incentiveRepo.Create(&incentive)
				if err != nil {
					return fmt.Errorf("error while storing incentive, err: %s", err.Error())
				}
			}

		}

		// update sync with last synced height
		sync.Incentives = lastBlock
		if err := sync.Save(); err != nil {
			return err
		}
	}

	fmt.Printf("incentives synced to block %d", sync.Incentives)

	return nil
}
