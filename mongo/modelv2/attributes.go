package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attribute struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	EventID      primitive.ObjectID `json:"event_id" bson:"event_id" validate:"required"`
	Key          string             `json:"key" bson:"key" validate:"required"`
	CompositeKey string             `json:"composite_key" bson:"composite_key" validate:"required"`
	Value        string             `json:"value" bson:"value" validate:"required"`
}

func (a *Attribute) Validate() error {
	return utility.ValidateStruct(&a)
}
