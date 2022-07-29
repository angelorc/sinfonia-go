package types

type Coin struct {
	Amount string `json:"amount" bson:"amount" validate:"required"`
	Denom  string `json:"denom" bson:"denom" validate:"required"`
}
