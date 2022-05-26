package model

import (
	"context"
	"errors"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/angelorc/sinfonia-go/mongo/db"
)

/**
 * DB Info
 */

const DB_COLLECTION_NAME__MESSAGE = "messages"
const DB_REF_NAME__MESSAGE = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__MESSAGE = []string{"from", "to"}

/**
 * MODEL
 */

type Message struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    int64              `json:"height" bson:"height"`
	TxHash    string             `json:"tx_hash" bson:"tx_hash"`
	MsgIndex  int                `json:"msg_index" bson:"msg_index"`
	MsgType   string             `json:"msg_type" bson:"msg_type"`
	Signer    string             `json:"signer" bson:"signer"`
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type MessageOrderByENUM string

/**
 * DTO
 */

// Read

type MessageWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type MessageWhere struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    int64               `json:"height,omitempty" bson:"height"`
	TxHash    string              `json:"tx_hash,omitempty" bson:"tx_hash"`
	MsgIndex  int                 `json:"msg_index,omitempty" bson:"msg_index"`
	MsgType   string              `json:"msg_type,omitempty" bson:"msg_type,omitempty"`
	Signer    string              `json:"signer,omitempty" bson:"signer,omitempty"`
	Timestamp time.Time           `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	OR        []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type MessageCreate struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    *int64              `json:"height,omitempty" bson:"height"`
	TxHash    *string             `json:"tx_hash" bson:"tx_hash" validate:"required"`
	MsgIndex  *int                `json:"msg_index" bson:"msg_index" validate:"required"`
	MsgType   *string             `json:"msg_type,omitempty" bson:"msg_type,omitempty"`
	Signer    *string             `json:"signer,omitempty" bson:"signer,omitempty"`
	Timestamp time.Time           `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Message) Message(filter *MessageWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__MESSAGE, DB_REF_NAME__MESSAGE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *Message) List(filter *MessageWhere, orderBy *MessageOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Message, error) {
	var items []*Message
	orderByKey := "timestamp"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__MESSAGE, DB_REF_NAME__MESSAGE)
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

func (m *Message) Count(filter *MessageWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__MESSAGE, DB_REF_NAME__MESSAGE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *Message) Create(data *MessageCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MESSAGE, DB_REF_NAME__MESSAGE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: check unique
	item := new(Transaction)
	f := bson.M{
		"$and": []bson.M{
			{"height": data.Height},
			{"tx_hash": data.TxHash},
			{"msg_index": data.MsgIndex},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Hash != "" {
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
