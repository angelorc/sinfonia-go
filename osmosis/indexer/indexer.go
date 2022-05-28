package indexer

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/avast/retry-go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v7/x/gamm/pool-models/balancer"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/angelorc/sinfonia-go/osmosis/chain"
	tmtypes "github.com/tendermint/tendermint/types"

	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
	incentivetypes "github.com/osmosis-labs/osmosis/v7/x/incentives/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 400)
	RtyErr    = retry.LastErrorOnly(true)
)

type IndexModules struct {
	Blocks       bool
	Transactions bool
	Messages     bool
	BlockResults bool
}

type Indexer struct {
	client     *chain.Client
	modules    *IndexModules
	concurrent int
}

func NewIndexer(client *chain.Client, modules *IndexModules, concurrent int) *Indexer {
	return &Indexer{
		client:     client,
		modules:    modules,
		concurrent: concurrent,
	}
}

func (i *Indexer) Parse(fromBlock, toBlock int64) {
	diff := (fromBlock - toBlock) + 1
	blocks := make([]int64, diff)
	for i := fromBlock; i <= toBlock; i++ {
		blocks[i-1] = i
	}

	if i.modules.Blocks {
		if err := i.parseBlocks(blocks, i.concurrent, i.IndexTransactions); err != nil {
			log.Fatalf("failed to index blocks. err: %v", err)
		}
	}
}

func (i *Indexer) parseBlocks(blocks []int64, concurrent int, cb func(height int64) error) error {
	fmt.Println("starting block queries for", i.client.ChainID())

	var (
		eg           errgroup.Group
		mutex        sync.Mutex
		failedBlocks = make([]int64, 0)
		sem          = make(chan struct{}, concurrent)
	)

	for _, height := range blocks {
		height := height
		sem <- struct{}{}

		eg.Go(func() error {
			if err := cb(height); err != nil {
				if strings.Contains(err.Error(), "wrong ID: no ID") {
					mutex.Lock()
					failedBlocks = append(failedBlocks, height)
					mutex.Unlock()
				} else {
					return fmt.Errorf("[height %d] - failed to get block. err: %s", height, err.Error())
				}
			}

			<-sem
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if len(failedBlocks) > 0 {
		return i.parseBlocks(failedBlocks, concurrent, cb)
	}

	return nil
}

func (i *Indexer) IndexTransactions(height int64) error {
	fmt.Println(fmt.Sprintf("Index transactions on block %d", height))

	block, err := i.client.QueryBlock(context.Background(), &height)
	if err != nil {
		if err = retry.Do(func() error {
			block, err = i.client.QueryBlock(context.Background(), &height)
			if err != nil {
				return err
			}

			return nil
		}, RtyAtt, RtyDel, RtyErr, retry.DelayType(retry.BackOffDelay), retry.OnRetry(func(n uint, err error) {
			log.Fatalf("retry: attempt %d, height %d, err: %v", n, height, err)
		})); err != nil {
			return err
		}
	}

	if block != nil {
		if i.modules.Transactions {
			i.parseTxs(block.Block.Height, block.Block.Data.Txs, block.Block.Time)
		}
	}

	if i.modules.BlockResults {
		if err := i.parseBlockResults(height, block.Block.Time); err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) parseBlockResults(height int64, ts time.Time) error {
	blockResults, err := i.client.QueryBlockResults(context.Background(), &height)
	if err != nil {
		if err = retry.Do(func() error {
			blockResults, err = i.client.QueryBlockResults(context.Background(), &height)
			if err != nil {
				return err
			}

			return nil
		}, RtyAtt, RtyDel, RtyErr, retry.DelayType(retry.BackOffDelay), retry.OnRetry(func(n uint, err error) {
			log.Fatalf("retry block_results: attempt %d, height %d, err: %v", n, height, err)
		})); err != nil {
			return err
		}
	}

	i.HandleBeginBlockEvents(height, blockResults.BeginBlockEvents, ts)

	return nil
}

func (i *Indexer) parseTxs(height int64, txs tmtypes.Txs, ts time.Time) {
	for index, tx := range txs {
		txTx, sdkTxRes, err := i.client.QueryTx(context.Background(), tx.Hash())
		if err != nil {
			log.Fatalf("[Height %d] {%d/%d txs} - Failed to query tx results. Err: %s \n", height, index+1, len(txs), err.Error())
		}

		feeAmt, feeDenom := i.client.ParseTxFee(txTx.GetFee())

		logStr := ""

		if sdkTxRes.Code > 0 {
			logStr = sdkTxRes.Logs.String()
		}

		err = i.InsertTx(tx.Hash(), sdkTxRes.Code, logStr, feeAmt, feeDenom, height, sdkTxRes.GasUsed, sdkTxRes.GasWanted, ts)
		if err != nil {
			log.Fatalf("[Height %d] {%d/%d txs} - Failed to write tx to db. Err: %s", height, index+1, len(txs), err.Error())
		}

		log.Printf("[Height %d] {%d/%d txs} - Successfuly wrote tx to db with %d msgs.", height, index+1, len(txs), len(txTx.GetMsgs()))

		if sdkTxRes.Code == 0 {
			for msgIndex, msg := range txTx.GetMsgs() {
				i.HandleMsg(msg, msgIndex, height, tx.Hash(), ts)
			}

			i.HandleLogs(sdkTxRes.Logs, txTx.GetMsgs(), height, tx.Hash(), ts)

			// i.HandleEvents(txTx.GetMsgs(), sdkTxRes.Events, height, tx.Hash(), ts)
		}
	}
}

func (i *Indexer) HandleMsg(msg sdk.Msg, msgIndex int, height int64, hash []byte, timestamp time.Time) {
	signer := i.client.MustEncodeAccAddr(msg.GetSigners()[0])

	err := i.InsertAccount(signer, timestamp)
	if err != nil {
		log.Fatalf("Failed to insert Account - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}

	err = i.InsertMsg(height, hash, msgIndex, sdk.MsgTypeURL(msg), signer, timestamp)
	if err != nil {
		log.Fatalf("Failed to insert MsgSend - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}

	/*switch m := msg.(type) {
	case *banktypes.MsgSend:
		err := i.InsertMsgSend(m.FromAddress, m.ToAddress, m.Amount.String(), hash, msgIndex, timestamp)
		if err != nil {
			log.Fatalf("Failed to insert MsgSend - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
		}

	case *stakingtypes.MsgCreateValidator:
	case *stakingtypes.MsgEditValidator:
	case *stakingtypes.MsgDelegate:
		err := i.InsertMsgDelegate(m.ValidatorAddress, m.DelegatorAddress, m.Amount.String(), hash, msgIndex, timestamp)
		if err != nil {
			log.Fatalf("Failed to insert MsgDelegate - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
		}

	case *stakingtypes.MsgUndelegate:
	case *stakingtypes.MsgBeginRedelegate:

	case *distrtypes.MsgSetWithdrawAddress:
	case *distrtypes.MsgWithdrawDelegatorReward:
	case *distrtypes.MsgWithdrawValidatorCommission:
	case *distrtypes.MsgFundCommunityPool:

	case *fantokentypes.MsgIssueFanToken:
		err := i.InsertMsgIssueFantoken(m.Name, m.Owner, m.MaxSupply.String(), m.Symbol, m.URI, hash, msgIndex, timestamp)
		if err != nil {
			log.Fatalf("Failed to insert MsgIssueFantoken - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
		}

	case *fantokentypes.MsgEditFanToken:
	case *fantokentypes.MsgMintFanToken:
	case *fantokentypes.MsgBurnFanToken:
	case *fantokentypes.MsgTransferFanTokenOwner:

	case *ibcclienttypes.MsgCreateClient:
	case *ibcclienttypes.MsgUpdateClient:
	case *ibcclienttypes.MsgSubmitMisbehaviour:
	case *ibcclienttypes.MsgUpgradeClient:

	case *ibcconnectiontypes.MsgConnectionOpenInit:
	case *ibcconnectiontypes.MsgConnectionOpenConfirm:
	case *ibcconnectiontypes.MsgConnectionOpenAck:
	case *ibcconnectiontypes.MsgConnectionOpenTry:

	case *channeltypes.MsgChannelOpenInit:
	case *channeltypes.MsgChannelOpenTry:
	case *channeltypes.MsgChannelOpenAck:
	case *channeltypes.MsgChannelOpenConfirm:
	case *channeltypes.MsgChannelCloseInit:
	case *channeltypes.MsgChannelCloseConfirm:
	case *channeltypes.MsgRecvPacket:
	case *channeltypes.MsgTimeout:
	case *channeltypes.MsgAcknowledgement:

	case *transfertypes.MsgTransfer:

	default:
		log.Fatalf("Unknown msg type: %s", sdk.MsgTypeURL(msg))
	}*/

}

func (i *Indexer) HandleLogs(logs sdk.ABCIMessageLogs, msgs []sdk.Msg, height int64, hash []byte, timestamp time.Time) {
	for index, mlog := range logs {
		i.HandleEvents(mlog.Events, msgs[index], index, height, hash, timestamp)
	}
}

func (i *Indexer) HandleEvents(events sdk.StringEvents, msg sdk.Msg, msgIndex int, height int64, hash []byte, ts time.Time) {
	for _, evt := range events {
		switch evt.Type {
		case gammtypes.TypeEvtTokenSwapped:
			i.handleTokenSwapped(height, hash, msgIndex, msg, evt.Attributes, ts)
		case gammtypes.TypeEvtPoolCreated:
			i.handlePoolCreated(height, hash, msgIndex, msg, evt.Attributes, ts)
		}
	}
}

func (i *Indexer) handleTokenSwapped(height int64, hash []byte, msgIndex int, _ sdk.Msg, attrs []sdk.Attribute, ts time.Time) {
	poolId := uint64(0)
	tokensIn := ""
	tokensOut := ""
	sender := ""

	for _, a := range attrs {
		switch a.Key {
		case sdk.AttributeKeySender:
			sender = a.Value
		case gammtypes.AttributeKeyPoolId:
			poolId, _ = strconv.ParseUint(a.Value, 0, 64)
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

func calcFee(tokenInStr, swapFeeStr string) string {
	tokenIn, _ := sdk.ParseCoinNormalized(tokenInStr)
	swapFee, _ := sdk.NewDecFromStr(swapFeeStr)
	tokenInAfterFee := tokenIn.Amount.ToDec().Mul(sdk.OneDec().Sub(swapFee)).TruncateInt()

	return sdk.Coin{
		Denom:  tokenIn.Denom,
		Amount: tokenIn.Amount.Sub(tokenInAfterFee),
	}.String()
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
}

func (i *Indexer) HandleBeginBlockEvents(height int64, events []abci.Event, ts time.Time) {
	for _, evt := range events {
		switch evt.Type {
		case incentivetypes.TypeEvtDistribution:
			i.handleIncentives(height, evt.Attributes, ts)
		}
	}
}

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
