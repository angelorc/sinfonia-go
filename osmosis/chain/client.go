package chain

import (
	"context"
	"encoding/hex"
	"fmt"
	tmcli "github.com/angelorc/sinfonia-go/tendermint"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v7/app"
	appparams "github.com/osmosis-labs/osmosis/v7/app/params"
)

type Client struct {
	config *Config
	rpc    rpcclient.Client
	codec  appparams.EncodingConfig
}

func NewClient(config *Config) (*Client, error) {
	timeout, _ := time.ParseDuration(config.Timeout)

	rpcClient, err := tmcli.NewClient(config.RPCAddr, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		config: config,
		codec:  app.MakeEncodingConfig(),
		rpc:    rpcClient,
	}, nil
}

func (c *Client) ChainID() string {
	return c.config.ChainID
}

func (c *Client) QueryBlock(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	return c.rpc.Block(ctx, height)
}

func (c *Client) QueryTx(ctx context.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return c.rpc.Tx(ctx, hash, prove)
}

func (c *Client) QueryTxFromString(ctx context.Context, hashHex string, prove bool) (*ctypes.ResultTx, error) {
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return nil, err
	}

	return c.QueryTx(ctx, hash, prove)
}

func (c *Client) DecodeTx(tx []byte) (sdk.Tx, error) {
	sdkTx, err := c.codec.TxConfig.TxDecoder()(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx. Err: %s \n", err.Error())
	}

	return sdkTx, nil
}

func (c *Client) DecodeTxFee(sdkTx sdk.Tx) (string, string) {
	fee := sdkTx.(sdk.FeeTx)

	var feeAmount, feeDenom string

	if len(fee.GetFee()) == 0 {
		feeAmount = "0"
		feeDenom = ""
	} else {
		feeAmount = fee.GetFee()[0].Amount.String()
		feeDenom = fee.GetFee()[0].Denom
	}

	return feeAmount, feeDenom
}
