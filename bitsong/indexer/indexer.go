package indexer

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/sync/errgroup"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/angelorc/sinfonia-go/bitsong/chain"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 400)
	RtyErr    = retry.LastErrorOnly(true)
)

type IndexModules struct {
	Blocks           bool
	Transactions     bool
	Messages         bool
	BeginBlockEvents bool
}

type Indexer struct {
	client  *chain.Client
	modules *IndexModules
}

func NewIndexer(client *chain.Client) *Indexer {
	return &Indexer{client: client}
}

func (i *Indexer) Parse(fromBlock, toBlock int64, concurrent int, modules *IndexModules) {
	var blocks []int64
	for i := fromBlock; i <= toBlock; i++ {
		blocks = append(blocks, i)
	}

	i.modules = modules

	if i.modules.Blocks {
		if err := i.parseBlocks(blocks, concurrent, i.IndexTransactions); err != nil {
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

	return nil
}

func (i *Indexer) parseTxs(height int64, txs tmtypes.Txs, ts time.Time) {
	for index, tx := range txs {
		sdkTx, err := i.client.DecodeTx(tx)
		if err != nil {
			log.Fatalf("[Height %d] {%d/%d txs} - %s", height, index+1, len(txs), err.Error())
		}

		txRes, err := i.client.QueryTx(context.Background(), tx.Hash(), false)
		if err != nil {
			log.Fatalf("[Height %d] {%d/%d txs} - Failed to query tx results. Err: %s \n", height, index+1, len(txs), err.Error())
		}

		feeAmt, feeDenom := i.client.DecodeTxFee(sdkTx)

		if txRes.TxResult.Code > 0 {
			logStr := fmt.Sprintf("{\"error\":\"%s\"}", txRes.TxResult.Log)

			err = i.InsertTx(tx.Hash(), logStr, feeAmt, feeDenom, height, txRes.TxResult.GasUsed, txRes.TxResult.GasWanted, ts, txRes.TxResult.Code)
			if err != nil {
				log.Fatalf("[Height %d] {%d/%d txs} - Failed to write tx to db. Err: %s", height, index+1, len(txs), err.Error())
			}
		} else {
			err = i.InsertTx(tx.Hash(), txRes.TxResult.Log, feeAmt, feeDenom, height, txRes.TxResult.GasUsed, txRes.TxResult.GasWanted, ts, txRes.TxResult.Code)

			if err != nil {
				log.Fatalf("[Height %d] {%d/%d txs} - Failed to write tx to db. Err: %s", height, index+1, len(txs), err.Error())
			}

			log.Printf("[Height %d] {%d/%d txs} - Successfuly wrote tx to db with %d msgs.", height, index+1, len(txs), len(sdkTx.GetMsgs()))
		}

		if i.modules.Messages {
			for msgIndex, msg := range sdkTx.GetMsgs() {
				i.HandleMsg(msg, msgIndex, height, tx.Hash(), ts)
			}
		}
	}
}

func (i *Indexer) HandleMsg(msg sdk.Msg, msgIndex int, height int64, hash []byte, timestamp time.Time) {
	fmt.Println(fmt.Sprintf("handle msg: %d, height %d", msgIndex, height))

	signer := i.client.MustEncodeAccAddr(msg.GetSigners()[0])

	err := i.InsertMsg(height, hash, msgIndex, sdk.MsgTypeURL(msg), signer, timestamp)
	if err != nil {
		log.Fatalf("Failed to insert MsgSend - index (%d), height (%d), err: %s", msgIndex, height, err.Error())
	}

}
