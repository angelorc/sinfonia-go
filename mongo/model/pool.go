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

const DB_COLLECTION_NAME__POOL = "pools"
const DB_REF_NAME__POOL = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__POOL = []string{"pool_id"}

/**
 * MODEL
 */

type Pool struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height   int64              `json:"height" bson:"height"`
	TxHash   string             `json:"tx_hash" bson:"tx_hash"`
	MsgIndex int                `json:"msg_index" bson:"msg_index"`

	PoolID     uint64      `json:"pool_id" bson:"pool_id" validate:"required"`
	PoolAssets []PoolAsset `json:"pool_assets" bson:"pool_assets" validate:"required"`
	SwapFee    string      `json:"swap_fee" bson:"swap_fee" validate:"required"`
	ExitFee    string      `json:"exit_fee" bson:"exit_fee"`
	Sender     string      `json:"sender" bson:"sender" validate:"required"`
	Timestamp  time.Time   `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

type PoolAsset struct {
	Token  string `json:"token" bson:"token" validate:"required"`
	Weight string `json:"weight" bson:"weight" validate:"required"`
}

/**
 * ENUM
 */

type PoolOrderByENUM string

/**
 * DTO
 */

// Read

type PoolWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type PoolWhere struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height   *int64              `json:"height,omitempty" bson:"height,omitempty"`
	TxHash   *string             `json:"tx_hash,omitempty" bson:"tx_hash,omitempty"`
	MsgIndex *int                `json:"msg_index,omitempty" bson:"msg_index,omitempty"`

	PoolID     *uint64      `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
	PoolAssets *[]PoolAsset `json:"pool_assets,omitempty" bson:"pool_assets,omitempty"`
	OR         []bson.M     `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type PoolCreate struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height   *int64              `json:"height,omitempty" bson:"height"`
	TxHash   *string             `json:"tx_hash" bson:"tx_hash" validate:"required"`
	MsgIndex *int                `json:"msg_index" bson:"msg_index" validate:"required"`

	PoolID     uint64      `json:"pool_id" bson:"pool_id" validate:"required"`
	PoolAssets []PoolAsset `json:"pool_assets" bson:"pool_assets" validate:"required"`
	SwapFee    string      `json:"swap_fee" bson:"swap_fee" validate:"required"`
	ExitFee    string      `json:"exit_fee" bson:"exit_fee"`
	Sender     string      `json:"sender" bson:"sender" validate:"required"`
	Timestamp  time.Time   `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Pool) One(filter *PoolWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__POOL, DB_REF_NAME__POOL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *Pool) List(filter *PoolWhere, orderBy *PoolOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Pool, error) {
	var items []*Pool
	orderByKey := "pool_id"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__POOL, DB_REF_NAME__POOL)
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

func (m *Pool) Count(filter *PoolWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__POOL, DB_REF_NAME__POOL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *Pool) Create(data *PoolCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__POOL, DB_REF_NAME__POOL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: check unique
	item := new(Pool)
	f := bson.M{
		"$and": []bson.M{
			{"pool_id": data.PoolID},
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
