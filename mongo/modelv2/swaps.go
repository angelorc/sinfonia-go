package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Swap struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	TxHash  string             `json:"tx_hash" bson:"tx_hash" validate:"required"`

	Account  string  `json:"account" bson:"account" validate:"required"`
	PoolId   int64   `json:"pool_id" bson:"pool_id" validate:"required"`
	Type     int     `json:"type" bson:"type"` // 0 - buy, 1 - sell
	TokenIn  Coin    `json:"token_in" bson:"token_in" validate:"required"`
	TokenOut Coin    `json:"token_out" bson:"token_out"`
	Fee      float64 `json:"fee" bson:"fee"`
	UsdValue float64 `json:"usd_value" bson:"usd_value"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (e *Swap) Validate() error {
	return utility.ValidateStruct(&e)
}

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

	Account  string  `json:"account" bson:"account" validate:"required"`
	PoolId   int64   `json:"pool_id" bson:"pool_id" validate:"required"`
	Type     int     `json:"type" bson:"type"` // 0 - buy, 1 - sell
	TokenIn  Coin    `json:"token_in" bson:"token_in" validate:"required"`
	TokenOut Coin    `json:"token_out" bson:"token_out" validate:"required"`
	Fee      float64 `json:"fee" bson:"fee"`
	UsdValue float64 `json:"usd_value" bson:"usd_value"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (ec *SwapCreateReq) Validate() error {
	return utility.ValidateStruct(ec)
}
