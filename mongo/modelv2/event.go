package modelv2

import (
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	MsgIndex   int         `json:"msg_index" bson:"msg_index"`
	Type       string      `json:"type" bson:"type"`
	Attributes []Attribute `json:"attributes" bson:"attributes" validate:"required"`
}

func (e *Event) Validate() error {
	return utility.ValidateStruct(&e)
}

type Attribute struct {
	Key   string `json:"key" bson:"key" validate:"required"`
	Value string `json:"value" bson:"value" validate:"required"`
}

func (a *Attribute) Validate() error {
	return utility.ValidateStruct(&a)
}

type ABCIMessageLog struct {
	MsgIndex int           `json:"msg_index,omitempty" bson:"msg_index,omitempty"`
	Log      string        `json:"log,omitempty" bson:"log,omitempty"`
	Events   []StringEvent `json:"events" bson:"events"`
}

type StringEvent struct {
	Type       string      `json:"type,omitempty" bson:"type,omitempty"`
	Attributes []Attribute `json:"attributes" bson:"attributes"`
}

type EventFilter struct {
	Id       *primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	TxID     *primitive.ObjectID `json:"tx_id,omitempty" bson:"tx_id,omitempty"`
	Key      *string             `json:"key,omitempty" bson:"key,omitempty"`
	MsgIndex *int                `json:"msg_index,omitempty" bson:"msg_index,omitempty"`
	Type     *string             `json:"type,omitempty" bson:"type,omitempty"`
}

func (ef *EventFilter) Validate() error {
	return nil
}

type EventCreateReq struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	TxID       primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex   int                `json:"msg_index" bson:"msg_index"`
	Type       string             `json:"type" bson:"type"`
	Attributes []Attribute        `json:"attributes" bson:"attributes" validate:"required"`
}

func (ec *EventCreateReq) Validate() error {
	return utility.ValidateStruct(ec)
}
