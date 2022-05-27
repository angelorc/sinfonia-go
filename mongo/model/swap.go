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
	Height   int64              `json:"height" bson:"height"`
	TxHash   string             `json:"tx_hash" bson:"tx_hash"`
	MsgIndex int                `json:"msg_index" bson:"msg_index"`

	PoolId    uint64 `json:"pool_id" bson:"pool_id"`
	TokensIn  string `json:"tokens_in" bson:"tokens_in"`
	TokensOut string `json:"tokens_out" bson:"tokens_out"`
	Account   string `json:"account" bson:"account"`
	Fee       string `json:"fee" bson:"fee"`

	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
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
	Height   int64               `json:"height,omitempty" bson:"height"`
	TxHash   string              `json:"tx_hash,omitempty" bson:"tx_hash"`
	MsgIndex int                 `json:"msg_index,omitempty" bson:"msg_index"`

	PoolId    uint64 `json:"pool_id" bson:"pool_id"`
	TokensIn  string `json:"tokens_in" bson:"tokens_in"`
	TokensOut string `json:"tokens_out" bson:"tokens_out"`
	Account   string `json:"account" bson:"account"`
	Fee       string `json:"fee" bson:"fee"`

	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	OR        []bson.M  `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type SwapCreate struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height   *int64              `json:"height,omitempty" bson:"height"`
	TxHash   *string             `json:"tx_hash" bson:"tx_hash" validate:"required"`
	MsgIndex *int                `json:"msg_index" bson:"msg_index" validate:"required"`

	PoolId    *uint64 `json:"pool_id" bson:"pool_id"`
	TokensIn  *string `json:"tokens_in" bson:"tokens_in"`
	TokensOut *string `json:"tokens_out" bson:"tokens_out"`
	Account   *string `json:"account" bson:"account"`
	Fee       *string `json:"fee" bson:"fee"`

	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Swap) Swap(filter *SwapWhere) error {
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

	// TODO: check unique
	item := new(Swap)
	f := bson.M{
		"$and": []bson.M{
			{"height": data.Height},
			{"tx_hash": data.TxHash},
			{"msg_index": data.MsgIndex},
			{"pool_id": data.PoolId},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.TxHash != "" {
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
