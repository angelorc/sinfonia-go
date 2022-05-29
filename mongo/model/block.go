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

const DB_COLLECTION_NAME__BLOCK = "blocks"
const DB_REF_NAME__BLOCK = "default"

/**
 * MODEL
 */

type Block struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64              `json:"height" bson:"height" validate:"required"`
	Hash    string             `json:"hash" bson:"hash" validate:"required"`
	Time    time.Time          `json:"time" bson:"time" validate:"required"`
}

/**
 * ENUM
 */

type BlockOrderByENUM string

/**
 * DTO
 */

// Read

type BlockWhere struct {
	ID      *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID *string             `json:"chain_id,omitempty" bson:"chain_id,omitempty"`
	Height  *int64              `json:"height,omitempty" bson:"height,omitempty"`
	Hash    *string             `json:"hash,omitempty" bson:"hash,omitempty"`
	Time    *time.Time          `json:"time,omitempty" bson:"time,omitempty"`
	OR      *[]bson.M           `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type BlockCreate struct {
	ID      *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID string              `json:"chain_id" bson:"chain_id" validate:"required"`
	Height  int64               `json:"height" bson:"height" validate:"required"`
	Hash    string              `json:"hash" bson:"hash" validate:"required"`
	Time    time.Time           `json:"time" bson:"time" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (b *Block) One(filter *BlockWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__BLOCK, DB_REF_NAME__BLOCK)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&b)

	return nil
}
func (b *Block) List(filter *BlockWhere, orderBy *BlockOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Block, error) {
	var items []*Block
	orderByKey := "height"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__BLOCK, DB_REF_NAME__BLOCK)
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
func (b *Block) Count(filter *BlockWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__BLOCK, DB_REF_NAME__BLOCK)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
func GetLastHeight() int64 {
	collection := db.GetCollection(DB_COLLECTION_NAME__BLOCK, DB_REF_NAME__BLOCK)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	block := new(Block)
	filter := &BlockWhere{}
	opts := options.FindOne().SetSort(bson.M{"height": -1})

	collection.FindOne(ctx, filter, opts).Decode(&block)
	return block.Height
}

// Write Operations

func (b *Block) Create(data *BlockCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__BLOCK, DB_REF_NAME__BLOCK)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check unique
	item := new(Block)
	f := bson.M{
		"$and": []bson.M{
			{"_id": data.ID},
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

/**
 * INDEXER API
 */

func InsertBlock(data *BlockCreate) error {
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := new(Block).Create(data); err != nil {
		return err
	}

	return nil
}
