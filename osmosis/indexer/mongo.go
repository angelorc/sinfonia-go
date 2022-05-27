package indexer

import (
	"encoding/hex"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/utility"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func (i *Indexer) InsertTx(hash []byte, code uint32, log, feeAmount, feeDenom string, height, gasUsed, gasWanted int64, timestamp time.Time) error {
	id, err := primitive.ObjectIDFromHex(hex.EncodeToString(hash[:12]))
	if err != nil {
		return err
	}

	hashStr := hex.EncodeToString(hash)

	item := model.Transaction{}
	data := model.TransactionCreate{
		ID:        &id,
		Height:    height,
		Hash:      &hashStr,
		Code:      code,
		Log:       &log,
		FeeAmount: &feeAmount,
		FeeDenom:  &feeDenom,
		GasUsed:   gasUsed,
		GasWanted: gasWanted,
		Timestamp: timestamp,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) InsertMsg(height int64, txHash []byte, msgIndex int, msgType, signer string, timestamp time.Time) error {
	hashStr := hex.EncodeToString(txHash)

	item := model.Message{}
	data := model.MessageCreate{
		Height:    &height,
		TxHash:    &hashStr,
		MsgIndex:  &msgIndex,
		MsgType:   &msgType,
		Signer:    &signer,
		Timestamp: timestamp,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) InsertAccount(acc string, firstSeen time.Time) error {
	item := model.Account{}
	data := model.AccountCreate{
		Address:   acc,
		FirstSeen: firstSeen,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) InsertSwap(height int64, txHash []byte, msgIndex int, poolId uint64, tokensIn, tokensOut, fee, acc string, ts time.Time) error {
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
