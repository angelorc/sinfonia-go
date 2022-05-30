package chain

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	tmcli "github.com/angelorc/sinfonia-go/tendermint"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"regexp"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/app"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"

	indexertypes "github.com/angelorc/sinfonia-go/indexer/types"
)

var _ indexertypes.ClientI = &Client{}

type Client struct {
	config *Config
	rpc    rpcclient.Client
	grpc   *grpc.ClientConn
	txSC   tx.ServiceClient
	codec  appparams.EncodingConfig
}

func NewClient(config *Config) (*Client, error) {
	timeout, _ := time.ParseDuration(config.Timeout)

	rpcClient, err := tmcli.NewClient(config.RPCAddr, timeout)
	if err != nil {
		return nil, err
	}

	// create grpc conn
	var grpcOpts []grpc.DialOption
	if config.GRPCInsecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	address := regexp.MustCompile("https?://").ReplaceAllString(config.GRPCAddr, "")
	grpcConn, err := grpc.Dial(address, grpcOpts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		config: config,
		codec:  app.MakeEncodingConfig(),
		rpc:    rpcClient,
		grpc:   grpcConn,
		txSC:   tx.NewServiceClient(grpcConn),
	}, nil
}

func (c *Client) ChainID() string {
	return c.config.ChainID
}

func (c *Client) QueryBlock(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	return c.rpc.Block(ctx, height)
}

func (c *Client) QueryBlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error) {
	return c.rpc.BlockResults(ctx, height)
}

func (c *Client) QueryTx(ctx context.Context, hash []byte) (*tx.Tx, *sdk.TxResponse, error) {
	res, err := c.txSC.GetTx(ctx, &tx.GetTxRequest{Hash: hex.EncodeToString(hash)})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tx. Err: %s \n", err.Error())
	}

	for _, msg := range res.Tx.Body.Messages {
		var stdMsg sdk.Msg
		err = c.codec.Marshaler.UnpackAny(msg, &stdMsg)
		if err != nil {
			return nil, nil, fmt.Errorf("error while unpacking message: %s", err)
		}
	}

	return res.Tx, res.TxResponse, nil
}

func (c *Client) QueryTxFromString(ctx context.Context, hashHex string) (*tx.Tx, *sdk.TxResponse, error) {
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return nil, nil, err
	}

	return c.QueryTx(ctx, hash)
}

func (c *Client) DecodeTx(tx []byte) (sdk.Tx, error) {
	sdkTx, err := c.codec.TxConfig.TxDecoder()(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx. Err: %s \n", err.Error())
	}

	return sdkTx, nil
}

/*func (c *Client) ParseTxFee(fees sdk.Coins) (string, string) {
	var feeAmount, feeDenom string

	if len(fees) == 0 {
		feeAmount = "0"
		feeDenom = ""
	} else {
		feeAmount = fees[0].Amount.String()
		feeDenom = fees[0].Denom
	}

	return feeAmount, feeDenom
}*/
