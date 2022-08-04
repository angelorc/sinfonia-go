package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LiquidityEvent struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	TxHash  string             `json:"tx_hash" bson:"tx_hash" validate:"required"`

	Type      string `json:"type" bson:"type"`
	Sender    string `json:"sender" bson:"sender" validate:"required"`
	PoolID    uint64 `json:"pool_id" bson:"pool_id" validate:"required"`
	TokensIn  []Coin `json:"tokens_in" bson:"tokens_in"`
	TokensOut []Coin `json:"tokens_out" bson:"tokens_out"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (e *LiquidityEvent) Validate() error {
	return utility.ValidateStruct(&e)
}

type LiquidityEventFilter struct {
	Id     *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Height *int64              `json:"height,omitempty" bson:"height,omitempty"`
	Sender *string             `json:"sender,omitempty" bson:"sender,omitempty"`
}

func (ef *LiquidityEventFilter) Validate() error {
	return nil
}

type LiquidityEventCreateReq struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	TxHash  string             `json:"tx_hash" bson:"tx_hash"`

	Sender    string `json:"sender" bson:"sender" validate:"required"`
	PoolID    uint64 `json:"pool_id" bson:"pool_id" validate:"required"`
	TokensIn  []Coin `json:"tokens_in" bson:"tokens_in"`
	TokensOut []Coin `json:"tokens_out" bson:"tokens_out"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (ec *LiquidityEventCreateReq) Validate() error {
	return utility.ValidateStruct(ec)
}
