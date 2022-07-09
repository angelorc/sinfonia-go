package model

import (
	"context"
	"errors"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

/**
 * DB Info
 */

const DB_COLLECTION_NAME__FANTOKEN = "fantokens"
const DB_REF_NAME__FANTOKEN = "default"

/**
 * MODEL
 */

type Fantoken struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   int64              `json:"height" bson:"height" validate:"required"`
	TxID     primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	Denom    string             `json:"denom" bson:"denom" validate:"required"`
	Alias    []string           `json:"alias" bson:"alias"`
	Owner    string             `json:"owner" bson:"owner" validate:"required"`
	IssuedAt time.Time          `json:"issued_at,omitempty" bson:"issued_at,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type FantokenOrderByENUM string

/**
 * DTO
 */

// Read

type FantokenWhere struct {
	ID      *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID *string             `json:"chain_id,omitempty" bson:"chain_id,omitempty"`
	Height  *int64              `json:"height,omitempty" bson:"height,omitempty"`
	TxID    *primitive.ObjectID `json:"tx_id,omitempty" bson:"tx_id,omitempty"`
	Denom   *string             `json:"denom,omitempty" bson:"denom,omitempty"`
	Alias   *string             `json:"alias,omitempty" bson:"alias,omitempty"`
	Owner   *string             `json:"owner,omitempty" bson:"owner,omitempty"`
	OR      []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type FantokenCreate struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChainID  *string             `json:"chain_id" bson:"chain_id" validate:"required"`
	Height   *int64              `json:"height" bson:"height" validate:"required"`
	TxID     *primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
	Denom    *string             `json:"denom" bson:"denom" validate:"required"`
	Alias    *[]string           `json:"alias" bson:"alias"`
	Owner    *string             `json:"owner" bson:"owner" validate:"required"`
	IssuedAt *time.Time          `json:"issued_at,omitempty" bson:"issued_at,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (f *Fantoken) One(filter *FantokenWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__FANTOKEN, DB_REF_NAME__FANTOKEN)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&f)

	return nil
}

func (f *Fantoken) List(filter *FantokenWhere, orderBy *FantokenOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Fantoken, error) {
	var items []*Fantoken
	orderByKey := "issued_at"
	orderByValue := -1

	collection := db.GetCollection(DB_COLLECTION_NAME__FANTOKEN, DB_REF_NAME__FANTOKEN)
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

func (f *Fantoken) Count(filter *FantokenWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__FANTOKEN, DB_REF_NAME__FANTOKEN)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Write Operations

func (f *Fantoken) Create(data *FantokenCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__FANTOKEN, DB_REF_NAME__FANTOKEN)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check unique
	item := new(FantokenCreate)
	filter := bson.M{
		"$and": []bson.M{
			{"denom": data.Denom},
		},
	}
	collection.FindOne(ctx, filter).Decode(&item)
	if item.Denom != nil {
		return nil
	}

	// operation
	if data.Alias == nil {
		data.Alias = &[]string{}
	}
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

func (f *Fantoken) AddAlias(alias string) error {
	alias = strings.TrimSpace(alias)

	// validate
	if utility.IsZeroVal(f.ID) {
		return errors.New("missing fantoken id")
	}

	if alias == "" {
		return errors.New("missing fantoken alias")
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__FANTOKEN, DB_REF_NAME__FANTOKEN)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check if fantoken exists
	collection.FindOne(ctx, bson.M{"_id": f.ID}).Decode(&f)
	if f.Denom == "" {
		return errors.New("fantoken not found")
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": f.ID}, bson.D{
		{"$push", bson.D{{"alias", alias}}},
	})
	if err != nil {
		return err
	}

	collection.FindOne(ctx, bson.M{"_id": f.ID}).Decode(&f)

	return nil
}
