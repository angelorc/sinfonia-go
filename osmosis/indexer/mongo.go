package indexer

import (
	"encoding/hex"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/utility"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"time"
)

func (i *Indexer) InsertSwap(height int64, txHash []byte, msgIndex int, poolId int64, tokensIn, tokensOut, fee, acc string, ts time.Time) error {
	hashStr := hex.EncodeToString(txHash)

	item := model.Swap{}
	data := model.SwapCreate{
		Height:    &height,
		TxHash:    &hashStr,
		MsgIndex:  &msgIndex,
		PoolId:    &poolId,
		TokensIn:  &tokensIn,
		TokensOut: &tokensOut,
		Fee:       &fee,
		Account:   &acc,
		Timestamp: ts,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) InsertPool(height int64, txHash []byte, msgIndex int, poolId uint64, poolAssets []model.PoolAsset, swapFee, exitFee, sender string, ts time.Time) error {
	hashStr := hex.EncodeToString(txHash)

	item := model.Pool{}
	data := model.PoolCreate{
		Height:     &height,
		TxHash:     &hashStr,
		MsgIndex:   &msgIndex,
		PoolID:     poolId,
		PoolAssets: poolAssets,
		SwapFee:    swapFee,
		ExitFee:    exitFee,
		Sender:     sender,
		Timestamp:  ts,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) InsertIncentive(height int64, receiver, coinsStr string, ts time.Time) error {
	coins, err := sdk.ParseCoinsNormalized(coinsStr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	assets := make([]model.IncentiveAsset, len(coins))
	for i, coin := range coins {
		assets[i].Amount = coin.Amount.Uint64()
		assets[i].Denom = coin.Denom
	}

	item := model.Incentive{}
	data := model.IncentiveCreate{
		Height:    height,
		Receiver:  receiver,
		Assets:    assets,
		Timestamp: ts,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}
