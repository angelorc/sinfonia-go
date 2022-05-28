package indexer

import (
	"encoding/hex"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (i *Indexer) InsertTx(hash []byte, log, feeAmount, feeDenom string, height, gasUsed, gasWanted int64, timestamp time.Time, code uint32) error {
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
		Height:    height,
		TxHash:    hashStr,
		MsgIndex:  msgIndex,
		MsgType:   msgType,
		Signer:    signer,
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
