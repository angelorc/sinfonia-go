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

const DB_COLLECTION_NAME__SWAP = "swaps"
const DB_REF_NAME__SWAP = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__SWAP = []string{"pool_id", "tokens_in", "tokens_out"}

/**
 * MODEL
 */

type Swap struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   int64              `json:"height" bson:"height" validate:"required"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex int                `json:"msg_index" bson:"msg_index"`

	PoolId    int64   `json:"pool_id" bson:"pool_id"`
	TokensIn  string  `json:"tokens_in" bson:"tokens_in"`
	TokensOut string  `json:"tokens_out" bson:"tokens_out"`
	Account   string  `json:"account" bson:"account"`
	Fee       string  `json:"fee" bson:"fee"`
	Volume    float64 `json:"volume" bson:"volume"`

	Time time.Time `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type SwapOrderByENUM string

/**
 * DTO
 */

// Read

type SwapWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type SwapWhere struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   *int64              `json:"height" bson:"height" validate:"required"`
	TxID     *primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex *int                `json:"msg_index,omitempty" bson:"msg_index,omitempty"`

	PoolId    *int64  `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
	TokensIn  *string `json:"tokens_in,omitempty" bson:"tokens_in,omitempty"`
	TokensOut *string `json:"tokens_out,omitempty" bson:"tokens_out,omitempty"`
	Account   *string `json:"account,omitempty" bson:"account,omitempty"`
	Fee       *string `json:"fee,omitempty" bson:"fee,omitempty"`

	OR []bson.M `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type SwapCreate struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   *int64              `json:"height" bson:"height" validate:"required"`
	TxID     *primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex *int                `json:"msg_index" bson:"msg_index" validate:"required"`

	PoolId    *int64   `json:"pool_id" bson:"pool_id"`
	TokensIn  *string  `json:"tokens_in" bson:"tokens_in"`
	TokensOut *string  `json:"tokens_out" bson:"tokens_out"`
	Account   *string  `json:"account" bson:"account"`
	Fee       *string  `json:"fee" bson:"fee"`
	Volume    *float64 `json:"volume" bson:"volume"`

	Time time.Time `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Swap) One(filter *SwapWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__SWAP, DB_REF_NAME__SWAP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *Swap) List(filter *SwapWhere, orderBy *SwapOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Swap, error) {
	var items []*Swap
	orderByKey := "timestamp"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__SWAP, DB_REF_NAME__SWAP)
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

func (m *Swap) Count(filter *SwapWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__SWAP, DB_REF_NAME__SWAP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *Swap) Create(data *SwapCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__SWAP, DB_REF_NAME__SWAP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: checkPools unique
	item := new(Swap)
	f := bson.M{
		"$and": []bson.M{
			{"height": data.Height},
			{"tx_id": data.TxID},
			{"msg_index": data.MsgIndex},
			{"pool_id": data.PoolId},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.ChainID != "" {
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
