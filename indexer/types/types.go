package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type ClientI interface {
	ChainID() string
	QueryBlock(ctx context.Context, height *int64) (*coretypes.ResultBlock, error)
	QueryBlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error)
	QueryTx(ctx context.Context, hash []byte) (*tx.Tx, *sdk.TxResponse, error)
	QueryTxFromString(ctx context.Context, hashHex string) (*tx.Tx, *sdk.TxResponse, error)
	EncodeBech32AccAddr(addr sdk.AccAddress) (string, error)
	MustEncodeAccAddr(addr sdk.AccAddress) string
	// ParseTxFee(fees sdk.Coins) (string, string)
	// DecodeTx(tx []byte) (sdk.Tx, error)
}
