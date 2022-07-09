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

const DB_COLLECTION_NAME__MERKLEDROP = "merkledrops"
const DB_REF_NAME__MERKLEDROP = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__MERKLEDROP = []string{"merkledrop_id"}

/**
 * MODEL
 */

type Merkledrop struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   int64              `json:"height" bson:"height" validate:"required"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex int                `json:"msg_index" bson:"msg_index"`

	MerkledropID int64  `json:"merkledrop_id" bson:"merkledrop_id"`
	Denom        string `json:"denom" bson:"denom"`
	Amount       int64  `json:"amount" bson:"amount"`
	StartHeight  int64  `json:"start_height" bson:"start_height"`
	EndHeight    int64  `json:"end_height" bson:"end_height"`

	Name  string `json:"name" bson:"name"`
	Image string `json:"image" bson:"image"`

	Time time.Time `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type MerkledropOrderByENUM string

/**
 * DTO
 */

// Read

type MerkledropWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type MerkledropWhere struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id" bson:"chain_id,omitempty"`
	Height   *int64              `json:"height" bson:"height,omitempty"`
	TxID     *primitive.ObjectID `json:"tx_id" bson:"tx_id,omitempty"`
	MsgIndex *int                `json:"msg_index,omitempty" bson:"msg_index,omitempty"`

	MerkledropID *int64  `json:"merkledrop_id,omitempty" bson:"merkledrop_id,omitempty"`
	Denom        *string `json:"denom,omitempty" bson:"denom,omitempty"`
	Amount       *int64  `json:"amount,omitempty" bson:"amount,omitempty"`
	StartHeight  *int64  `json:"start_height,omitempty" bson:"start_height,omitempty"`
	EndHeight    *int64  `json:"end_height,omitempty" bson:"end_height,omitempty"`

	Name  *string `json:"name,omitempty" bson:"name,omitempty"`
	Image *string `json:"image,omitempty" bson:"image,omitempty"`

	OR []bson.M `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type MerkledropCreate struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   *int64              `json:"height" bson:"height" validate:"required"`
	TxID     *primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex *int                `json:"msg_index" bson:"msg_index" validate:"required"`

	MerkledropId *int64  `json:"merkledrop_id" bson:"merkledrop_id" validate:"required"`
	Denom        *string `json:"denom" bson:"denom" validate:"required"`
	Amount       *int64  `json:"amount" bson:"amount" validate:"required"`
	StartHeight  *int64  `json:"start_height" bson:"start_height" validate:"required"`
	EndHeight    *int64  `json:"end_height" bson:"end_height" validate:"required"`
	Name         *string `json:"name,omitempty" bson:"name,omitempty"`
	Image        *string `json:"image,omitempty" bson:"image,omitempty"`

	Time *time.Time `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Merkledrop) One(filter *MerkledropWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP, DB_REF_NAME__MERKLEDROP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *Merkledrop) List(filter *MerkledropWhere, orderBy *MerkledropOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Merkledrop, error) {
	var items []*Merkledrop
	orderByKey := "time"
	orderByValue := -1

	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP, DB_REF_NAME__MERKLEDROP)
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

func (m *Merkledrop) Count(filter *MerkledropWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP, DB_REF_NAME__MERKLEDROP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *Merkledrop) Create(data *MerkledropCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP, DB_REF_NAME__MERKLEDROP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: checkPools unique
	item := new(Swap)
	f := bson.M{
		"$and": []bson.M{
			{"height": data.Height},
			{"tx_id": data.TxID},
			{"msg_index": data.MsgIndex},
			{"merkledrop_id": data.MerkledropId},
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
