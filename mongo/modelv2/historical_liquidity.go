package modelv2

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type HistoricalLiquidity struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PoolID uint64             `json:"pool_id" bson:"pool_id"`
	Assets []Coin             `json:"assets" bson:"assets" validate:"required"`
	Time   time.Time          `json:"time" bson:"time" validate:"required"`
}

type HistoricalLiquidityCreateReq struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PoolID uint64             `json:"pool_id" bson:"pool_id"`
	Assets []Coin             `json:"assets" bson:"assets" validate:"required"`
	Time   time.Time          `json:"time" bson:"time" validate:"required"`
}
