package model

type Coin struct {
	Amount string `json:"amount" bson:"amount" validate:"required"`
	Denom  string `json:"denom" bson:"denom" validate:"required"`
}

type CoinInput Coin
