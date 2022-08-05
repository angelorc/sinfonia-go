package modelv2

import (
	"fmt"
)

type Coin struct {
	Amount float64 `json:"amount" bson:"amount"`
	Denom  string  `json:"denom" bson:"denom"`
}

func (c Coin) String() string {
	return fmt.Sprintf("%f%s", c.Amount, c.Denom)
}
