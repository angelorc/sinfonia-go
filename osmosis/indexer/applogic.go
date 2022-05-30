package indexer

import (
	incentivetypes "github.com/osmosis-labs/osmosis/v7/x/incentives/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
	"time"
)

/* func (i *Indexer) handleTokenSwapped(height int64, hash []byte, msgIndex int, _ sdk.Msg, attrs []sdk.Attribute, ts time.Time) {
	poolId := int64(0)
	tokensIn := ""
	tokensOut := ""
	sender := ""

	for _, a := range attrs {
		switch a.Key {
		case sdk.AttributeKeySender:
			sender = a.Value
		case gammtypes.AttributeKeyPoolId:
			poolId, _ = strconv.ParseInt(a.Value, 0, 64)
		case gammtypes.AttributeKeyTokensIn:
			tokensIn = a.Value
		case gammtypes.AttributeKeyTokensOut:
			tokensOut = a.Value
		}
	}

	item := new(model.Pool)
	item.One(
		&model.PoolWhere{
			PoolID: &poolId,
		},
	)

	err := i.InsertSwap(height, hash, msgIndex, poolId, tokensIn, tokensOut, calcFee(tokensIn, item.SwapFee), sender, ts)
	if err != nil {
		log.Fatalf("Failed to insert TokenSwap - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}
}

func (i *Indexer) handlePoolCreated(height int64, hash []byte, msgIndex int, msg sdk.Msg, attrs []sdk.Attribute, ts time.Time) {
	msgCreateBalancerPool := msg.(*balancer.MsgCreateBalancerPool)
	poolId := uint64(0)
	for _, a := range attrs {
		switch a.Key {
		case gammtypes.AttributeKeyPoolId:
			poolId, _ = strconv.ParseUint(string(a.Value), 0, 64)
		}
	}

	poolAssets := make([]model.PoolAsset, len(msgCreateBalancerPool.PoolAssets))

	for i, pa := range msgCreateBalancerPool.PoolAssets {
		poolAsset := model.PoolAsset{
			Token:  pa.Token.String(),
			Weight: pa.Weight.String(),
		}
		poolAssets[i] = poolAsset
	}

	swapFee := msgCreateBalancerPool.PoolParams.SwapFee
	exitFee := msgCreateBalancerPool.PoolParams.ExitFee
	sender := i.client.MustEncodeAccAddr(msg.GetSigners()[0])

	err := i.InsertPool(height, hash, msgIndex, poolId, poolAssets, swapFee.String(), exitFee.String(), sender, ts)
	if err != nil {
		log.Fatalf("Failed to insert TokenSwap - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}
}*/

func (i *Indexer) handleIncentives(height int64, attrs []abci.EventAttribute, ts time.Time) {
	receiver := ""
	coins := ""

	for _, attr := range attrs {
		switch string(attr.Key) {
		case incentivetypes.AttributeReceiver:
			receiver = string(attr.Value)
		case incentivetypes.AttributeAmount:
			coins = string(attr.Value)
		}
	}

	err := i.InsertIncentive(height, receiver, coins, ts)
	if err != nil {
		log.Fatalf("Failed to insert Incentive - height (%d), receiver (%s), err: %s", height, receiver, err.Error())
	}
}
