package modelv2

import (
	"fmt"
	"strconv"
)

type Coin struct {
	Amount string `json:"amount" bson:"amount"`
	Denom  string `json:"denom" bson:"denom"`
}

func (c Coin) String() string {
	return fmt.Sprintf("%s%s", c.Amount, c.Denom)
}

func (c Coin) GetAmount() float64 {
	amt, _ := strconv.ParseFloat(c.Amount, 64)
	return amt
}
