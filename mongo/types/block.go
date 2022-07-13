package types

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockFilter struct {
	Id     *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Height *int64              `json:"height,omitempty" bson:"height,omitempty"`
}

func (bf *BlockFilter) Validate() error {
	return nil
}

type BlockCreateReq struct{}

func (bc *BlockCreateReq) Validate() error {
	return utility.ValidateStruct(&bc)
}
