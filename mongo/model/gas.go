package model

type Gas struct {
	Used   int64 `json:"used" bson:"used" validate:"required"`
	Wanted int64 `json:"wanted" bson:"wanted" validate:"required"`
}
