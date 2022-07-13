package indexer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	types2 "github.com/angelorc/sinfonia-go/mongo/types"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/angelorc/sinfonia-go/indexer/types"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/scalar"
	"github.com/avast/retry-go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
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
	client     types.ClientI
	modules    *IndexModules
	concurrent int
}

func NewIndexer(client types.ClientI, modules *IndexModules, concurrent int) *Indexer {
	return &Indexer{
		client:     client,
		modules:    modules,
		concurrent: concurrent,
	}
}

func (i *Indexer) Parse(fromBlock, toBlock int64) {
	blocks := make([]int64, 0)

	for i := fromBlock; i <= toBlock; i++ {
		blocks = append(blocks, i)
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

	if block == nil {
		return fmt.Errorf("block not found")
	}

	blockRepo := repository.NewBlockRepository()
	blockDB, err := blockRepo.Create(&types2.BlockCreateReq{
		ChainID: block.Block.ChainID,
		Height:  block.Block.Height,
		Hash:    block.BlockID.Hash.String(),
		Time:    block.Block.Time,
	})

	if err != nil {
		log.Fatalf("[Height %d] - Failed to write block to db. Err: %s", height, err.Error())
	}

	if i.modules.Transactions {
		i.parseTxs(blockDB.ID, block.Block.ChainID, height, block.Block.Data.Txs, block.Block.Time)
	}

	if i.modules.BlockResults {
		if err := i.parseBlockResults(height, block.Block.Time); err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) parseTxs(blockID primitive.ObjectID, chainID string, height int64, txs tmtypes.Txs, time time.Time) {
	for index, tx := range txs {
		txTx, sdkTxRes, err := i.client.QueryTx(context.Background(), tx.Hash())
		if err != nil {
			log.Fatalf("[Height %d] {%d/%d txs} - Failed to query tx results. Err: %s \n", height, index+1, len(txs), err.Error())
		}

		if sdkTxRes.Code > 0 {
			continue
		}

		txID := model.TxHashToObjectID(tx.Hash())
		hashStr := hex.EncodeToString(tx.Hash())

		data := &model.TransactionCreate{
			ID:        &txID,
			ChainID:   &chainID,
			Height:    height,
			BlockID:   &blockID,
			Hash:      &hashStr,
			Code:      sdkTxRes.Code,
			Logs:      ConvertABCIMessageLogs(sdkTxRes.Logs),
			Fee:       ConvertCoins(txTx.GetFee()),
			GasUsed:   &sdkTxRes.GasUsed,
			GasWanted: &sdkTxRes.GasWanted,
			Time:      time,
		}

		err = model.InsertTx(data)
		if err != nil {
			log.Fatalf("[Height %d] {%d/%d txs} - Failed to write tx to db. Err: %s", height, index+1, len(txs), err.Error())
		}

		log.Printf("[Height %d] {%d/%d txs} - Successfuly wrote tx to db with %d msgs.", height, index+1, len(txs), len(txTx.GetMsgs()))

		for msgIndex, msg := range txTx.GetMsgs() {
			i.HandleMsg(txID, chainID, msg, msgIndex, height, time)
		}

		i.HandleLogs(sdkTxRes.Logs, txTx.GetMsgs(), height, tx.Hash(), time)
	}
}

func (i *Indexer) HandleMsg(txID primitive.ObjectID, chainID string, msg sdk.Msg, msgIndex int, height int64, time time.Time) {
	signer := i.client.MustEncodeAccAddr(msg.GetSigners()[0])

	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	msgType := sdk.MsgTypeURL(msg)
	msgValueStr, err := marshaler.MarshalToString(msg)
	if err != nil {
		log.Fatalf("Failed to marshal msg - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}

	var msgValue scalar.JSON
	err = json.Unmarshal([]byte(msgValueStr), &msgValue)
	if err != nil {
		log.Fatalf("Failed to marshaljson msg  - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}

	data := &model.MessageCreate{
		TxID:     &txID,
		Height:   &height,
		ChainID:  &chainID,
		MsgIndex: &msgIndex,
		MsgType:  &msgType,
		MsgValue: &msgValue,
		Signer:   &signer,
		Time:     time,
	}
	err = model.InsertMsg(data)
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

func (i *Indexer) HandleEvents(events sdk.StringEvents, msg sdk.Msg, msgIndex int, height int64, hash []byte, ts time.Time) {
	for _, evt := range events {
		switch evt.Type {
		//case gammtypes.TypeEvtTokenSwapped:
		// i.handleTokenSwapped(height, hash, msgIndex, msg, evt.Attributes, ts)
		//case gammtypes.TypeEvtPoolCreated:
		// i.handlePoolCreated(height, hash, msgIndex, msg, evt.Attributes, ts)
		}
	}
}

func (i *Indexer) HandleBeginBlockEvents(height int64, events []abci.Event, ts time.Time) {
	for _, evt := range events {
		switch evt.Type {
		//case incentivetypes.TypeEvtDistribution:
		// i.handleIncentives(height, evt.Attributes, ts)
		}
	}
}
