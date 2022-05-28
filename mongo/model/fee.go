package model

type Fee struct {
	Amount string `json:"amount" bson:"amount"`
	Denom  string `json:"denom" bson:"denom"`
}
