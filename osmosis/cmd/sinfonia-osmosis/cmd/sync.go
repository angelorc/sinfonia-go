package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/osmosis/chain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v9/x/gamm/pool-models/balancer"
	"github.com/osmosis-labs/osmosis/v9/x/gamm/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
		GetSyncPoolCmd(),
		GetSyncSwapCmd(),
	)

	return cmd
}

func GetSyncPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools",
		Short:   "sync pools from latest blocks",
		Example: "sinfonia-osmosis sync pools --mongo-dbname sinfonia-test",
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

func GetSyncSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swaps",
		Short:   "sync swaps from latest blocks",
		Example: "sinfonia-osmosis sync swaps --mongo-dbname sinfonia-test",
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

			if err := syncSwaps(); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func parseAttrs(attrs []model.Attribute) (poolID int64, tokensIn, tokensOut string) {
	for _, attr := range attrs {
		switch attr.Key {
		case "pool_id":
			poolID, _ = strconv.ParseInt(attr.Value, 0, 64)
		case "tokens_in":
			tokensIn = attr.Value
		case "tokens_out":
			tokensOut = attr.Value
		}
	}
	return poolID, tokensIn, tokensOut
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

func syncSwaps() error {
	// get last available height on db
	lastBlock := model.GetLastHeight()

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Swaps = int64(0)
	}

	txsLogs, err := model.GetTxsAndLogsByMessageType("/osmosis.gamm.v1beta1.MsgSwapExactAmountIn", sync.Swaps, lastBlock)
	if err != nil {
		return err
	}

	for _, txLogs := range txsLogs {
		for _, txlog := range txLogs.Tx.Logs {
			for _, evt := range txlog.Events {
				switch evt.Type {
				case "token_swapped":
					poolID, tokensIn, tokensOut := parseAttrs(evt.Attributes)

					pool := new(model.Pool)
					pool.One(
						&model.PoolWhere{
							PoolID: &poolID,
						},
					)

					fee := calcFee(tokensIn, pool.SwapFee)
					volume := calcVolumeUSD(tokensIn, tokensOut, txLogs.Time)

					// get volume by day
					// {
					//    _id: {
					//      date: {$dateToString: { format: "%Y-%m-%d", date: "$time" }},
					//    },
					//    volume: {$sum: "$volume"}
					// }
					// {
					//   "_id.date": 1
					// }

					// get volume by pool
					// {
					//    _id: {
					//      pool: "$pool_id",
					//    },
					//    volume: {$sum: "$volume"}
					// }
					// {
					//   "_id.volume": -1
					// }

					swapModel := new(model.Swap)
					data := &model.SwapCreate{
						ChainID:   &txLogs.ChainID,
						Height:    &txLogs.Height,
						TxID:      &txLogs.TxID,
						MsgIndex:  &txLogs.MsgIndex,
						PoolId:    &poolID,
						TokensIn:  &tokensIn,
						TokensOut: &tokensOut,
						Account:   &txLogs.Signer,
						Fee:       &fee,
						Volume:    &volume,
						Time:      txLogs.Time,
					}

					if err := swapModel.Create(data); err != nil {
						return err
					}
				}
			}
		}
	}

	// update sync with last synced height
	sync.Swaps = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("%d swaps synced to block %d ", len(txsLogs), sync.Swaps)

	return nil
}

func calcFee(tokenInStr, swapFeeStr string) string {
	tokenIn, _ := sdk.ParseCoinNormalized(tokenInStr)
	swapFee, _ := sdk.NewDecFromStr(swapFeeStr)
	tokenInAfterFee := tokenIn.Amount.ToDec().Mul(sdk.OneDec().Sub(swapFee)).TruncateInt()

	return sdk.Coin{
		Denom:  tokenIn.Denom,
		Amount: tokenIn.Amount.Sub(tokenInAfterFee),
	}.String()
}

func syncPools(client *chain.Client) error {
	// get last available height on db
	lastBlock := model.GetLastHeight()

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Pools = int64(0)
	}

	txsLogs, err := model.GetTxsAndLogsByMessageType("/osmosis.gamm.poolmodels.balancer.v1beta1.MsgCreateBalancerPool", sync.Pools, lastBlock)
	if err != nil {
		return err
	}

	for _, txLogs := range txsLogs {
		for _, txlog := range txLogs.Tx.Logs {
			for _, evt := range txlog.Events {
				switch evt.Type {
				case "pool_created":
					poolID, err := strconv.ParseUint(evt.Attributes[0].Value, 0, 64)
					if err != nil {
						return fmt.Errorf("error while parsing poolID, err: %s", err.Error())
					}

					poolRes, err := client.QueryPoolByID(poolID)
					if err != nil {
						return fmt.Errorf("error while fetching poolID, err: %s", err.Error())
					}

					var poolI types.PoolI
					err = client.Codec.Marshaler.UnpackAny(poolRes.GetPool(), &poolI)
					if err != nil {
						log.Fatalf("error while decoding the new pool")
					}

					pool, ok := poolI.(*balancer.Pool)
					if !ok {
						log.Fatalf("error while decoding the new pool")
					}

					allPoolAssets := pool.GetAllPoolAssets()

					// check if is contains a bitsong fantoken
					found := false
					for _, asset := range allPoolAssets {
						// if is an ibc token
						if strings.HasPrefix(asset.Token.Denom, "ibc/") {
							// check fantoken by alias
							var fantoken model.Fantoken

							if err := fantoken.One(&model.FantokenWhere{
								Alias: &asset.Token.Denom,
							}); err != nil {
								log.Fatalf("error while querying fantoken: %v", err)
							}

							// if exist return true
							if !fantoken.ID.IsZero() {
								found = true
							} else {
								// else  query ibc denom
								denomRes, err := client.QueryIBCDenomTrace(strings.ReplaceAll(asset.Token.Denom, "ibc/", ""))
								if err != nil {
									return fmt.Errorf("error while fetching denom-trace, err: %s", err.Error())
								}

								if strings.HasPrefix(denomRes.DenomTrace.BaseDenom, "ft") {
									if strings.HasPrefix(denomRes.DenomTrace.Path, "transfer/channel-0") {
										if err := fantoken.One(&model.FantokenWhere{
											Denom: &denomRes.DenomTrace.BaseDenom,
										}); err != nil {
											log.Fatalf("error while querying fantoken: %v", err)
										}

										// add the alias
										if err := fantoken.AddAlias(asset.Token.Denom); err != nil {
											log.Fatalf("error while adding a fantoken alias: %v", err)
										}

										found = true
									}
								}
							}
						}
					}

					if !found {
						continue
					}

					poolModel := new(model.Pool)
					data := &model.PoolCreate{
						ChainID:    &txLogs.ChainID,
						Height:     &txLogs.Height,
						TxID:       &txLogs.TxID,
						MsgIndex:   &txLogs.MsgIndex,
						PoolID:     poolID,
						PoolAssets: convertPoolAssetsToModel(allPoolAssets),
						SwapFee:    pool.GetSwapFee(sdk.Context{}).String(),
						ExitFee:    pool.GetExitFee(sdk.Context{}).String(),
						Sender:     txLogs.Signer,
						Time:       txLogs.Time,
					}

					if err := poolModel.Create(data); err != nil {
						return err
					}
				}
			}
		}
	}

	// update sync with last synced height
	sync.Pools = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("%d pools synced to block %d ", len(txsLogs), sync.Pools)

	return nil
}

func convertPoolAssetsToModel(pa []balancer.PoolAsset) []model.PoolAsset {
	newPoolAssets := make([]model.PoolAsset, len(pa))

	for i, p := range pa {
		newPoolAssets[i] = model.PoolAsset{
			Token:  model.Coin{Denom: p.Token.Denom, Amount: p.Token.Amount.String()},
			Weight: p.Weight.String(),
		}
	}

	return newPoolAssets
}
