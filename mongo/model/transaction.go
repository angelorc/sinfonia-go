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

const DB_COLLECTION_NAME__TRANSACTION = "transactions"
const DB_REF_NAME__TRANSACTION = "default"

/**
 * MODEL
 */

type Transaction struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    int64              `json:"height,omitempty" bson:"height,omitempty" validate:"required"`
	Hash      string             `json:"hash" bson:"hash" validate:"required"`
	Code      int                `json:"code" bson:"code"  validate:"required"`
	Log       string             `json:"log" bson:"log" validate:"required"`
	FeeAmount string             `json:"fee_amount" bson:"fee_amount"`
	FeeDenom  string             `json:"fee_denom" bson:"fee_denom"`
	GasUsed   int64              `json:"gas_used" bson:"gas_used"`
	GasWanted int64              `json:"gas_wanted" bson:"gas_wanted"`
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type TransactionOrderByENUM string

/**
 * DTO
 */

// Read

type TransactionWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type TransactionWhere struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    int64               `json:"height,omitempty" bson:"height,omitempty"`
	Hash      *string             `json:"hash,omitempty" bson:"hash,omitempty"`
	Code      int                 `json:"code,omitempty" bson:"code,omitempty"`
	Log       *string             `json:"log,omitempty" bson:"log,omitempty"`
	FeeAmount *string             `json:"fee_amount,omitempty" bson:"fee_amount,omitempty"`
	FeeDenom  *string             `json:"fee_denom,omitempty" bson:"fee_denom,omitempty"`
	GasUsed   int64               `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64               `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Timestamp time.Time           `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	OR        []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type TransactionCreate struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    int64               `json:"height,omitempty" bson:"height,omitempty" validate:"required"`
	Hash      *string             `json:"hash" bson:"hash,omitempty" validate:"required"`
	Code      uint32              `json:"code,omitempty" bson:"code,omitempty"`
	Log       *string             `json:"log" bson:"log" validate:"required"`
	FeeAmount *string             `json:"fee_amount,omitempty" bson:"fee_amount,omitempty"`
	FeeDenom  *string             `json:"fee_denom,omitempty" bson:"fee_denom,omitempty"`
	GasUsed   int64               `json:"gas_used,omitempty" bson:"gas_used,omitempty"`
	GasWanted int64               `json:"gas_wanted,omitempty" bson:"gas_wanted,omitempty"`
	Timestamp time.Time           `json:"timestamp,omitempty" bson:"timestamp,omitempty" validate:"required"`
}

type TransactionUpdate struct {
	Hash      *string   `json:"hash" bson:"hash,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
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

func (t *Transaction) Update(where primitive.ObjectID, data *TransactionUpdate) error {
	// validate
	if utility.IsZeroVal(where) {
		return errors.New("internal server error")
	}
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check transaction is exists
	collection.FindOne(ctx, bson.M{"_id": where}).Decode(&t)
	if t.Hash == "" {
		return errors.New("item not found")
	}

	// check unique
	item := new(Transaction)
	f := bson.M{
		"$or": []bson.M{
			{"hash": data.Hash, "_id": bson.M{"$ne": where}},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Hash != "" {
		return errors.New("transaction already exist")
	}

	// operation
	_, err := collection.UpdateOne(ctx, bson.M{"_id": where}, bson.M{"$set": data})
	collection.FindOne(ctx, bson.M{"_id": where}).Decode(&t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Delete() error {
	collection := db.GetCollection(DB_COLLECTION_NAME__TRANSACTION, DB_REF_NAME__TRANSACTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if utility.IsZeroVal(t.ID) {
		return errors.New("invalid id")
	}

	collection.FindOne(ctx, bson.M{"_id": t.ID}).Decode(&t)
	if t.Hash == "" {
		return errors.New("item not found")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": t.ID})
	if err != nil {
		return err
	}

	return nil
}

/**
 * INDEXER API
 */

func (t *Transaction) InsertTx(hash []byte, log, feeAmount, feeDenom string, height, gasUsed, gasWanted int64, timestamp time.Time, code uint32) error {
	id, err := primitive.ObjectIDFromHex(hex.EncodeToString(hash[:12]))
	if err != nil {
		return err
	}

	hashStr := hex.EncodeToString(hash)

	item := Transaction{}
	data := TransactionCreate{
		ID:        &id,
		Height:    height,
		Hash:      &hashStr,
		Code:      code,
		Log:       &log,
		FeeAmount: &feeAmount,
		FeeDenom:  &feeDenom,
		GasUsed:   gasUsed,
		GasWanted: gasWanted,
		Timestamp: timestamp,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}
