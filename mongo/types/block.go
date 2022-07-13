package types

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BlockFilter struct {
	Id     *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Height *int64              `json:"height,omitempty" bson:"height,omitempty"`
}

func (bf *BlockFilter) Validate() error {
	return nil
}

type BlockCreateReq struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	Hash    string             `json:"hash" bson:"hash" validate:"required"`
	Time    time.Time          `json:"time" bson:"time" validate:"required"`
}

func (bc *BlockCreateReq) Validate() error {
	return utility.ValidateStruct(bc)
}
