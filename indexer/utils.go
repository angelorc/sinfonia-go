package indexer

import (
	"github.com/angelorc/sinfonia-go/mongo/model"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ConvertCoin(coin sdk.Coin) model.Coin {
	return model.Coin{
		Amount: coin.Amount.String(),
		Denom:  coin.Denom,
	}
}

func ConvertCoins(coins sdk.Coins) *[]model.Coin {
	newCoins := make([]model.Coin, len(coins))

	for i, coin := range coins {
		newCoins[i] = ConvertCoin(coin)
	}

	return &newCoins
}

func ConvertABCIMessageLogs(logs sdk.ABCIMessageLogs) []model.ABCIMessageLog {
	newLogs := make([]model.ABCIMessageLog, len(logs))

	for i, log := range logs {
		newLogs[i] = model.ABCIMessageLog{
			MsgIndex: int(log.MsgIndex),
			Log:      log.Log,
			Events:   ConvertStringEvents(log.Events),
		}
	}

	return newLogs
}

func ConvertStringEvents(events sdk.StringEvents) []model.StringEvent {
	newEvents := make([]model.StringEvent, len(events))

	for i, evt := range events {
		newEvents[i] = model.StringEvent{
			Type:       evt.Type,
			Attributes: ConvertAttributes(evt.Attributes),
		}
	}

	return newEvents
}

func ConvertAttributes(attrs []sdk.Attribute) []model.Attribute {
	newAttrs := make([]model.Attribute, len(attrs))

	for i, attr := range attrs {
		newAttrs[i] = model.Attribute{
			Key:   attr.Key,
			Value: attr.Value,
		}
	}

	return newAttrs
}
