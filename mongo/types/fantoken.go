package types

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Asset struct {
}

type FantokenFilter struct {
	Id    *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Denom *string             `json:"denom,omitempty" bson:"denom,omitempty"`
}

func (ff *FantokenFilter) Validate() error {
	return nil
}

type FantokenCreateReq struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	ChainInfoCreateReq
	Denom    string    `json:"denom" bson:"denom" validate:"required"`
	IssuedAt time.Time `json:"issued_at" bson:"issued_at" validate:"required"`
}

func (fc *FantokenCreateReq) Validate() error {
	return utility.ValidateStruct(fc)
}
