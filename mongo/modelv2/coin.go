package modelv2

import "fmt"

type Coin struct {
	Amount string `json:"amount" bson:"amount" validate:"required"`
	Denom  string `json:"denom" bson:"denom" validate:"required"`
}

func (c Coin) String() string {
	return fmt.Sprintf("%s%s", c.Amount, c.Denom)
}
