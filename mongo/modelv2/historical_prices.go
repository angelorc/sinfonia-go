package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Price struct {
	Usd string `json:"usd" bson:"usd" validate:"required"`
}

type HistoricalPrice struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Asset string             `json:"asset" bson:"asset" validate:"required"`
	Price []Price            `json:"price" bson:"price" validate:"required"`
	Time  time.Time          `json:"time" bson:"time" validate:"required"`
}

func (b *HistoricalPrice) Validate() error {
	return utility.ValidateStruct(&b)
}

type HistoricalPriceFilter struct {
	Id    *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Asset *string             `json:"asset,omitempty" bson:"asset,omitempty"`
	Time  *time.Time          `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
}

func (bf *HistoricalPriceFilter) Validate() error {
	return nil
}

type HistoricalPriceCreateReq struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Asset string             `json:"asset" bson:"asset" validate:"required"`
	Price []Price            `json:"price" bson:"price" validate:"required"`
	Time  time.Time          `json:"time" bson:"time" validate:"required"`
}

func (bc *HistoricalPriceCreateReq) Validate() error {
	return utility.ValidateStruct(bc)
}
