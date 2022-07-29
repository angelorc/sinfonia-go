package modelv2

import (
	"github.com/angelorc/sinfonia-go/mongo/types"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	Hash    string             `json:"hash" bson:"hash" validate:"required"`
	Code    int                `json:"code" bson:"code"`
	//Logs      []ABCIMessageLog   `json:"logs" bson:"logs" validate:"required"`
	Events    []Event      `json:"events" bson:"events"`
	Fee       []types.Coin `json:"fee" bson:"fee"`
	GasUsed   int64        `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64        `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Time      time.Time    `json:"time" bson:"time" validate:"required"`
}

func (b *Transaction) Validate() error {
	return utility.ValidateStruct(&b)
}

type TransactionEvents struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id"`
	Height  int64              `json:"height" bson:"height"`
	Hash    string             `json:"hash" bson:"hash"`
	Time    time.Time          `json:"time" bson:"time"`
	Events  []Event            `json:"events" bson:"events"`
}

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
	Events    []Event            `json:"events" bson:"events"`
	Fee       []Coin             `json:"fee" bson:"fee"`
	GasUsed   int64              `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64              `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Time      time.Time          `json:"time" bson:"time" validate:"required"`
}

func (t *TransactionCreateReq) Validate() error {
	return utility.ValidateStruct(t)
}
