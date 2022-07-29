package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Pool struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height"`
	TxHash  string             `json:"tx_hash" bson:"tx_hash"`

	PoolID     uint64      `json:"pool_id" bson:"pool_id" validate:"required"`
	PoolAssets []PoolAsset `json:"pool_assets" bson:"pool_assets" validate:"required"`
	SwapFee    string      `json:"swap_fee" bson:"swap_fee" validate:"required"`
	ExitFee    string      `json:"exit_fee" bson:"exit_fee"`

	Time time.Time `json:"time,omitempty" bson:"time,omitempty"`
}

func (e *Pool) Validate() error {
	return utility.ValidateStruct(&e)
}

type PoolAsset struct {
	Token  Coin   `json:"token" bson:"token" validate:"required"`
	Weight string `json:"weight" bson:"weight" validate:"required"`
}

type PoolFilter struct {
	Id     *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	PoolID *uint64             `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
}

func (ef *PoolFilter) Validate() error {
	return nil
}

type PoolCreateReq struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height"`
	TxHash  string             `json:"tx_hash" bson:"tx_hash"`

	PoolID     uint64      `json:"pool_id" bson:"pool_id" validate:"required"`
	PoolAssets []PoolAsset `json:"pool_assets" bson:"pool_assets" validate:"required"`
	SwapFee    string      `json:"swap_fee" bson:"swap_fee" validate:"required"`
	ExitFee    string      `json:"exit_fee" bson:"exit_fee"`

	Time time.Time `json:"time" bson:"time"`
}

func (ec *PoolCreateReq) Validate() error {
	return utility.ValidateStruct(ec)
}
