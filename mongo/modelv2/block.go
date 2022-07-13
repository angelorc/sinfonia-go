package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Block struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	Hash    string             `json:"hash" bson:"hash" validate:"required"`
	Time    time.Time          `json:"time" bson:"time" validate:"required"`
}

func (b *Block) Validate() error {
	return utility.ValidateStruct(&b)
}
