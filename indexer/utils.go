package indexer

import (
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"
)

func ConvertCoin(coin sdk.Coin) modelv2.Coin {
	return modelv2.Coin{
		Amount: coin.Amount.String(),
		Denom:  coin.Denom,
	}
}

func ConvertCoins(coins sdk.Coins) *[]modelv2.Coin {
	newCoins := make([]modelv2.Coin, len(coins))

	for i, coin := range coins {
		newCoins[i] = ConvertCoin(coin)
	}

	return &newCoins
}

func GetAttrByKey(key string, attrs []sdk.Attribute) *string {
	for _, attr := range attrs {
		if attr.Key == key {
			return &attr.Value
		}
	}

	return nil
}

func IsAllowedTx(allowedActions []string, logs sdk.ABCIMessageLogs) bool {
	for _, log := range logs {
		for _, evt := range log.Events {
			if evt.Type == "message" {
				action := GetAttrByKey("action", evt.Attributes)

				if action != nil {
					if slices.Contains(allowedActions, *action) {
						return true
					}
				}
			}
		}
	}

	return false
}

func ConvertABCIMessageLogs(logs sdk.ABCIMessageLogs) []modelv2.ABCIMessageLog {
	newLogs := make([]modelv2.ABCIMessageLog, len(logs))

	for i, log := range logs {
		newLogs[i] = modelv2.ABCIMessageLog{
			MsgIndex: int(log.MsgIndex),
			Log:      log.Log,
			Events:   ConvertStringEvents(log.Events),
		}
	}

	return newLogs
}

func ConvertStringEvents(events sdk.StringEvents) []modelv2.StringEvent {
	newEvents := make([]modelv2.StringEvent, len(events))

	for i, evt := range events {
		newEvents[i] = modelv2.StringEvent{
			Type:       evt.Type,
			Attributes: ConvertAttributes(evt.Attributes),
		}
	}

	return newEvents
}

func ConvertAttributes(attrs []sdk.Attribute) []modelv2.Attribute {
	newAttrs := make([]modelv2.Attribute, len(attrs))

	for i, attr := range attrs {
		newAttrs[i] = modelv2.Attribute{
			Key:   attr.Key,
			Value: attr.Value,
		}
	}

	return newAttrs
}
