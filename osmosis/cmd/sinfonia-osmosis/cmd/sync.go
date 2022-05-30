package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/osmosis/chain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v7/x/gamm/types"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
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
		Use:     "pool",
		Short:   "sync pools from latest blocks",
		Example: "sinfonia-osmosis sync pool --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, mongoDBName, mongoRetryWrites, err := parseMongoFlags(cmd)
			if err != nil {
				return err
			}

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

			client, err := chain.NewClient(chain.GetOsmosisConfig())
			if err != nil {
				return fmt.Errorf("failed to get RPC endpoints on chain %s. err: %v", "osmosis", err)
			}

			if err := syncPools(client); err != nil {
				return err
			}

			return nil
		},
	}

	addMongoFlags(cmd)

	return cmd
}

func GetSyncSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swaps",
		Short:   "sync swaps from latest blocks",
		Example: "sinfonia-osmosis sync swaps --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, mongoDBName, mongoRetryWrites, err := parseMongoFlags(cmd)
			if err != nil {
				return err
			}

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

			if err := syncSwaps(); err != nil {
				return err
			}

			return nil
		},
	}

	addMongoFlags(cmd)

	return cmd
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

	txsLogs, err := model.GetTxsAndLogsByMessageType("/osmosis.gamm.v1beta1.MsgCreateBalancerPool", sync.Pools, lastBlock)
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

					var pool types.PoolI
					err = client.Codec.Marshaler.UnpackAny(poolRes.GetPool(), &pool)
					if err != nil {
						log.Fatalf("error while decoding the new pool")
					}

					poolAssets := make([]model.PoolAsset, len(pool.GetAllPoolAssets()))

					for i, pa := range pool.GetAllPoolAssets() {
						poolAsset := model.PoolAsset{
							Token:  pa.Token.String(),
							Weight: pa.Weight.String(),
						}
						poolAssets[i] = poolAsset
					}

					poolModel := new(model.Pool)
					data := &model.PoolCreate{
						ChainID:    &txLogs.ChainID,
						Height:     &txLogs.Height,
						TxID:       &txLogs.TxID,
						MsgIndex:   &txLogs.MsgIndex,
						PoolID:     poolID,
						PoolAssets: convertPoolAssetsToModel(pool.GetAllPoolAssets()),
						SwapFee:    pool.GetPoolSwapFee().String(),
						ExitFee:    pool.GetPoolExitFee().String(),
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

func convertPoolAssetsToModel(pa []types.PoolAsset) []model.PoolAsset {
	newPoolAssets := make([]model.PoolAsset, len(pa))

	for i, p := range pa {
		newPoolAssets[i] = model.PoolAsset{
			Token:  p.Token.String(),
			Weight: p.Weight.String(),
		}
	}

	return newPoolAssets
}

func getAttrValueByKey(key string, attrs []model.Attribute) string {
	for _, attr := range attrs {
		if key == attr.Key {
			return attr.Value
		}
	}

	return ""
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

					poolID, err := strconv.ParseInt(getAttrValueByKey("pool_id", evt.Attributes), 0, 64)
					if err != nil {
						return fmt.Errorf("error while parsing poolID, err: %s", err.Error())
					}

					tokensIn := getAttrValueByKey("tokens_in", evt.Attributes)
					tokensOut := getAttrValueByKey("tokens_out", evt.Attributes)

					pool := new(model.Pool)
					pool.One(
						&model.PoolWhere{
							PoolID: &poolID,
						},
					)

					fee := calcFee(tokensIn, pool.SwapFee)

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
