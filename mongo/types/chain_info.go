package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChainInfoCreateReq struct {
	ChainID  string             `json:"chain_id" bson:"chain_id,omitempty" validate:"required"`
	Height   int64              `json:"height" bson:"height,omitempty" validate:"required"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id,omitempty" validate:"required"`
	MsgIndex int                `json:"msg_index,omitempty" bson:"msg_index,omitempty"`
}
