package model

import (
	"context"
	"errors"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

/**
 * DB Info
 */

const DB_COLLECTION_NAME__INCENTIVE = "incentives"
const DB_REF_NAME__INCENTIVE = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__INCENTIVE = []string{"receiver"}

/**
 * MODEL
 */

type Incentive struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height int64              `json:"height" bson:"height"`

	Receiver  string           `json:"receiver" bson:"receiver" validate:"required"`
	Assets    []IncentiveAsset `json:"assets" bson:"assets"`
	Timestamp time.Time        `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

type IncentiveAsset struct {
	Amount int64  `json:"amount" bson:"amount"`
	Denom  string `json:"denom" bson:"denom"`
}

type IncentiveAssetWhere struct {
	Amount *int64  `json:"amount,omitempty" bson:"amount,omitempty"`
	Denom  *string `json:"denom,omitempty" bson:"denom,omitempty"`
}

/**
 * ENUM
 */

type IncentiveOrderByENUM string

/**
 * DTO
 */

// Read

type IncentiveWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type IncentiveWhere struct {
	ID        *primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    *int64                 `json:"height,omitempty" bson:"height,omitempty"`
	Receiver  *string                `json:"receiver,omitempty" bson:"receiver,omitempty"`
	Assets    *[]IncentiveAssetWhere `json:"assets,omitempty" bson:"assets,omitempty"`
	Timestamp *time.Time             `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	OR        []bson.M               `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type IncentiveCreate struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    int64              `json:"height,omitempty" bson:"height" validate:"required"`
	Receiver  string             `json:"receiver" bson:"receiver" validate:"required"`
	Assets    []IncentiveAsset   `json:"assets" bson:"assets" validate:"required"`
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Incentive) One(filter *IncentiveWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__INCENTIVE, DB_REF_NAME__INCENTIVE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *Incentive) List(filter *IncentiveWhere, orderBy *IncentiveOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Incentive, error) {
	var items []*Incentive
	orderByKey := "timestamp"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__INCENTIVE, DB_REF_NAME__INCENTIVE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := options.Find()
	if limit != nil {
		options.SetLimit(int64(*limit))
	}
	if skip != nil {
		options.SetSkip(int64(*skip))
	}
	if orderBy != nil {
		orderByKey, orderByValue = utility.GetOrderByKeyAndValue(string(*orderBy))
	}
	options.SetSort(map[string]int{orderByKey: orderByValue})

	var queryFilter interface{}
	if filter != nil {
		queryFilter = filter
	}
	if !utility.IsZeroVal(customQuery) {
		queryFilter = customQuery
	}

	cursor, err := collection.Find(ctx, &queryFilter, options)
	if err != nil {
		return items, err
	}
	err = cursor.All(ctx, &items)
	if err != nil {
		return items, err
	}

	return items, nil
}

func (m *Incentive) Count(filter *IncentiveWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__INCENTIVE, DB_REF_NAME__INCENTIVE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *Incentive) Create(data *IncentiveCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__INCENTIVE, DB_REF_NAME__INCENTIVE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: check unique
	item := new(Pool)
	f := bson.M{
		"$and": []bson.M{
			{"receiver": data.Receiver},
			{"height": data.Height},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.PoolID > 0 {
		return nil
	}

	// operation
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("server error")
	}

	return nil
}
