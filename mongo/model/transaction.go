package model

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
 * DB Info
 */

const (
	DB_COLLECTION_NAME__TRANSACTION = "transactions"
	DB_REF_NAME__TRANSACTION        = "default"
)

/**
 * MODEL
 */

type Transaction struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BlockID primitive.ObjectID `json:"block_id" bson:"block_id" validate:"required"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height,omitempty" bson:"height,omitempty" validate:"required"`
	Hash    string             `json:"hash" bson:"hash" validate:"required"`
	Code    int                `json:"code" bson:"code"  validate:"required"`
	Log     interface{}        `json:"log" bson:"log" validate:"required"`
	Fee     Fee                `json:"fee" bson:"fee"`
	Gas     Gas                `json:"gas" bson:"gas"`
	Time    time.Time          `json:"time" bson:"time" validate:"required"`
}

/**
 * ENUM
 */

type TransactionOrderByENUM string

/**
 * DTO
 */

// Read

type TransactionWhere struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BlockID   *primitive.ObjectID `json:"block_id,omitempty" bson:"block_id,omitempty"`
	ChainID   *string             `json:"chain_id,omitempty" bson:"chain_id,omitempty"`
	Height    *int64              `json:"height,omitempty" bson:"height,omitempty"`
	Hash      *string             `json:"hash,omitempty" bson:"hash,omitempty"`
	Code      int                 `json:"code,omitempty" bson:"code,omitempty"`
	Log       *string             `json:"log,omitempty" bson:"log,omitempty"`
	FeeAmount *string             `json:"fee_amount,omitempty" bson:"fee_amount,omitempty"`
	FeeDenom  *string             `json:"fee_denom,omitempty" bson:"fee_denom,omitempty"`
	GasUsed   int64               `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64               `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Time      time.Time           `json:"time,omitempty" bson:"time,omitempty"`
	OR        []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type TransactionCreate struct {
	ID      *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BlockID *primitive.ObjectID `json:"block_id" bson:"block_id" validate:"required"`
	ChainID *string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64               `json:"height,omitempty" bson:"height,omitempty" validate:"required"`
	Hash    *string             `json:"hash" bson:"hash" validate:"required"`
	Code    uint32              `json:"code" bson:"code"`
	Log     interface{}         `json:"log" bson:"log" validate:"required"`
	Fee     *Fee                `json:"fee,omitempty" bson:"fee,omitempty"`
	Gas     *Gas                `json:"gas,omitempty" bson:"gas,omitempty"`
	Time    time.Time           `json:"time,omitempty" bson:"time,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (t *Transaction) One(filter *TransactionWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&t)

	return nil
}
func (t *Transaction) List(filter *TransactionWhere, orderBy *TransactionOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Transaction, error) {
	var items []*Transaction
	orderByKey := "timestamp"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
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
func (t *Transaction) Count(filter *TransactionWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (t *Transaction) Create(data *TransactionCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: check unique
	item := new(Transaction)
	f := bson.M{
		"$and": []bson.M{
			{"hash": data.Hash},
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

	// collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&t)
	return nil
}

/**
 * INDEXER API
 */

func TxHashToObjectID(hash []byte) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex.EncodeToString(hash[:12]))
	if err != nil {
		panic(err)
	}

	return id
}

func InsertTx(data *TransactionCreate) error {
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := new(Transaction).Create(data); err != nil {
		return err
	}

	return nil
}
