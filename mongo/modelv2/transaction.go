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
	Fee       []types.Coin `json:"fee" bson:"fee"`
	GasUsed   int64        `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64        `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Time      time.Time    `json:"time" bson:"time" validate:"required"`
}

func (b *Transaction) Validate() error {
	return utility.ValidateStruct(&b)
}
