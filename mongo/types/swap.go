package types

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SwapFilter struct {
	Id     *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Height *int64              `json:"height,omitempty" bson:"height,omitempty"`
}

func (ef *SwapFilter) Validate() error {
	return nil
}

type SwapCreateReq struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	TxHash  string             `json:"tx_hash" bson:"tx_hash" validate:"required"`

	PoolId         int64  `json:"pool_id" bson:"pool_id" validate:"required"`
	TokensInAmt    string `json:"tokens_in_amt" bson:"tokens_in_amt" validate:"required"`
	TokensInDenom  string `json:"tokens_in_denom" bson:"tokens_in_denom" validate:"required"`
	TokensOutAmt   string `json:"tokens_out_amt" bson:"tokens_out_amt" validate:"required"`
	TokensOutDenom string `json:"tokens_out_denom" bson:"tokens_out_denom" validate:"required"`
	Account        string `json:"account" bson:"account" validate:"required"`
	Fee            string `json:"fee" bson:"fee"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (ec *SwapCreateReq) Validate() error {
	return utility.ValidateStruct(ec)
}
