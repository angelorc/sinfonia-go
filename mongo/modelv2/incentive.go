package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Incentive struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`

	Receiver string `json:"receiver" bson:"receiver" validate:"required"`
	Assets   []Coin `json:"assets" bson:"assets"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (e *Incentive) Validate() error {
	return utility.ValidateStruct(&e)
}

type IncentiveFilter struct {
	Id       *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Height   *int64              `json:"height,omitempty" bson:"height,omitempty"`
	Receiver *string             `json:"receiver,omitempty" bson:"receiver,omitempty"`
}

func (ef *IncentiveFilter) Validate() error {
	return nil
}

type IncentiveCreateReq struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`

	Receiver string `json:"receiver" bson:"receiver" validate:"required"`
	Assets   []Coin `json:"assets" bson:"assets"`

	Time time.Time `json:"time" bson:"time" validate:"required"`
}

func (ec *IncentiveCreateReq) Validate() error {
	return utility.ValidateStruct(ec)
}
