package types

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TransactionFilter struct {
	Id   *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Hash *string             `json:"hash,omitempty" bson:"hash,omitempty"`
}

func (tf *TransactionFilter) Validate() error {
	return nil
}

type TransactionCreateReq struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID   string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height    int64              `json:"height" bson:"height" validate:"required"`
	Hash      string             `json:"hash" bson:"hash" validate:"required"`
	Code      int                `json:"code" bson:"code"`
	Fee       []Coin             `json:"fee" bson:"fee"`
	GasUsed   int64              `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64              `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Time      time.Time          `json:"time" bson:"time" validate:"required"`
}

func (t *TransactionCreateReq) Validate() error {
	return utility.ValidateStruct(t)
}
