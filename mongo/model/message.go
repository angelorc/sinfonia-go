package model

import (
	"context"
	"errors"
	"github.com/angelorc/sinfonia-go/server/scalar"
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
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   int64              `json:"height" bson:"height" validate:"required"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex *int               `json:"msg_index" bson:"msg_index"`
	MsgType  string             `json:"msg_type" bson:"msg_type"`
	MsgValue scalar.JSON        `json:"msg_value" bson:"msg_value"`
	Signer   string             `json:"signer" bson:"signer"`
	Time     time.Time          `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
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
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id,omitempty" bson:"chain_id,omitempty"`
	Height   *int64              `json:"height,omitempty" bson:"height,omitempty"`
	TxID     *primitive.ObjectID `json:"tx_id,omitempty" bson:"tx_id,omitempty"`
	MsgIndex *int                `json:"msg_index,omitempty" bson:"msg_index,omitempty"`
	MsgType  *string             `json:"msg_type,omitempty" bson:"msg_type,omitempty"`
	MsgValue *scalar.JSON        `json:"msg_value,omitempty" bson:"msg_value,omitempty"`
	Signer   *string             `json:"signer,omitempty" bson:"signer,omitempty"`
	Time     *time.Time          `json:"time,omitempty" bson:"time,omitempty"`
	OR       *[]bson.M           `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type MessageCreate struct {
	ID       primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   *int64              `json:"height" bson:"height" validate:"required"`
	TxID     *primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	MsgIndex *int                `json:"msg_index" bson:"msg_index" validate:"required"`
	MsgType  *string             `json:"msg_type" bson:"msg_type" validate:"required"`
	MsgValue *scalar.JSON        `json:"msg_value" bson:"msg_value" validate:"required"`
	Signer   *string             `json:"signer" bson:"signer" validate:"required"`
	Time     time.Time           `json:"time" bson:"time" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Message) One(filter *MessageWhere) error {
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
	item := new(Message)
	f := bson.M{
		"$and": []bson.M{
			{"tx_id": data.TxID},
			{"msg_index": data.MsgIndex},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.TxID.String() == data.TxID.String() {
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

/**
 * INDEXER API
 */

func InsertMsg(data *MessageCreate) error {
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := new(Message).Create(data); err != nil {
		return err
	}

	return nil
}

// TxLogs struct
type TxLogs struct {
	Signer   string             `bson:"signer"`
	Time     time.Time          `bson:"time"`
	ChainID  string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   int64              `json:"height" bson:"height" validate:"required"`
	MsgIndex int                `json:"msg_index" bson:"msg_index"`
	MsgType  string             `json:"msg_type" bson:"msg_type"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	Tx       struct {
		Logs []struct {
			Events []struct {
				Type string `bson:"type"`

				Attributes []Attribute `bson:"attributes"`
			} `bson:"events"`
		} `bson:"logs"`
	} `bson:"tx"`
}

func GetTxsAndLogsByMessageType(msgType string, fromBlock, toBlock int64) ([]TxLogs, error) {
	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MESSAGE, DB_REF_NAME__MESSAGE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"height": bson.M{
					"$gt":  fromBlock,
					"$lte": toBlock,
				},
			},
		},
		{
			"$match": bson.M{
				"msg_type": bson.M{
					"$eq": msgType,
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "transactions",
				"localField":   "tx_id",
				"foreignField": "_id",
				"as":           "tx",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$tx",
			},
		},
		{
			"$project": bson.M{
				"chain_id":       1,
				"height":         1,
				"tx_id":          1,
				"msg_index":      1,
				"msg_type":       1,
				"signer":         1,
				"time":           1,
				"tx.logs.events": 1,
			},
		},
	}

	var txsLogs []TxLogs

	// aggregate pipeline
	accCursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return txsLogs, err
	}

	// decode
	if err = accCursor.All(ctx, &txsLogs); err != nil {
		return txsLogs, err
	}

	return txsLogs, nil
}
