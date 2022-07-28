package types

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttributeFilter struct {
	Id      *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	EventID *primitive.ObjectID `json:"event_id,omitempty" bson:"event_id,omitempty"`
	Key     *string             `json:"key,omitempty" bson:"key,omitempty"`
}

func (af *AttributeFilter) Validate() error {
	return nil
}

type AttributeCreateReq struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	EventID      primitive.ObjectID `json:"event_id" bson:"event_id" validate:"required"`
	Key          string             `json:"key" bson:"key" validate:"required"`
	CompositeKey string             `json:"composite_key" bson:"composite_key" validate:"required"`
	Value        string             `json:"value" bson:"value" validate:"required"`
}

func (ac *AttributeCreateReq) Validate() error {
	return utility.ValidateStruct(ac)
}
