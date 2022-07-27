package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Fantoken struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   int64              `json:"height" bson:"height" validate:"required"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	Denom    string             `json:"denom" bson:"denom" validate:"required"`
	Alias    []string           `json:"alias" bson:"alias"`
	Owner    string             `json:"owner" bson:"owner" validate:"required"`
	IssuedAt time.Time          `json:"issued_at" bson:"issued_at" validate:"required"`
}

func (f *Fantoken) Validate() error {
	return utility.ValidateStruct(&f)
}
