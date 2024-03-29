package chain

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	tmcli "github.com/angelorc/sinfonia-go/tendermint"
	"github.com/cosmos/cosmos-sdk/types/tx"
	gammtypes "github.com/osmosis-labs/osmosis/v9/x/gamm/types"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"regexp"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v9/app"
	appparams "github.com/osmosis-labs/osmosis/v9/app/params"
	"google.golang.org/grpc"

	"github.com/angelorc/sinfonia-go/indexer/types"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	ibctypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

var _ types.ClientI = &Client{}

type Client struct {
	config *config.ChainConfig
	rpc    rpcclient.Client
	grpc   *grpc.ClientConn
	txSC   tx.ServiceClient
	Codec  appparams.EncodingConfig
}

func NewClient(config *config.ChainConfig) (*Client, error) {
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
		Codec:  app.MakeEncodingConfig(),
		rpc:    rpcClient,
		grpc:   grpcConn,
		txSC:   tx.NewServiceClient(grpcConn),
	}, nil
}

func (c *Client) ChainID() string {
	return c.config.ChainID
}

func (c *Client) LatestBlockHeight(ctx context.Context) int64 {
	status, err := c.rpc.Status(ctx)
	if err != nil {
		return 0
	}
	return status.SyncInfo.LatestBlockHeight
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
		err = c.Codec.Marshaler.UnpackAny(msg, &stdMsg)
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
	sdkTx, err := c.Codec.TxConfig.TxDecoder()(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx. Err: %s \n", err.Error())
	}

	return sdkTx, nil
}

// APP Query

func (c *Client) QueryPoolByID(poolID uint64) (*gammtypes.QueryPoolResponse, error) {
	return gammtypes.NewQueryClient(c.grpc).Pool(context.Background(), &gammtypes.QueryPoolRequest{PoolId: poolID})
}

func (c *Client) QueryPoolByIDWithHeight(poolID uint64, height int64) (*gammtypes.QueryPoolResponse, error) {
	return gammtypes.NewQueryClient(c.grpc).Pool(
		metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, fmt.Sprintf("%d", height)),
		&gammtypes.QueryPoolRequest{PoolId: poolID},
	)
}

func (c *Client) QueryIBCDenomTrace(hash string) (*ibctypes.QueryDenomTraceResponse, error) {
	return ibctypes.NewQueryClient(c.grpc).DenomTrace(context.Background(), &ibctypes.QueryDenomTraceRequest{Hash: hash})
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

func (c *Client) EncodeBech32AccAddr(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(c.config.AccountPrefix, addr)
}
func (c *Client) MustEncodeAccAddr(addr sdk.AccAddress) string {
	enc, err := c.EncodeBech32AccAddr(addr)
	if err != nil {
		panic(err)
	}
	return enc
}
