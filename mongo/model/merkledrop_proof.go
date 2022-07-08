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

const DB_COLLECTION_NAME__MERKLEDROP_PROOF = "merkledrop_proofs"
const DB_REF_NAME__MERKLEDROP_PROOF = "default"

/**
 * SEARCH regex fields
 */

var SEARCH_FILEDS__MERKLEDROP_PROOF = []string{"merkledrop_id"}

/**
 * MODEL
 */

type MerkledropProof struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	MerkledropID int64    `json:"merkledrop_id" bson:"merkledrop_id"`
	Index        int64    `json:"index" bson:"index"`
	Address      string   `json:"address" bson:"address"`
	Amount       int64    `json:"amount" bson:"amount"`
	Proofs       []string `json:"proofs" bson:"proofs"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type MerkledropProofOrderByENUM string

/**
 * DTO
 */

// Read

type MerkledropProofWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type MerkledropProofWhere struct {
	ID *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	MerkledropID *int64    `json:"merkledrop_id,omitempty" bson:"merkledrop_id,omitempty"`
	Address      *string   `json:"address,omitempty" bson:"address,omitempty"`
	Index        *int64    `json:"index,omitempty" bson:"index,omitempty"`
	Amount       *int64    `json:"amount,omitempty" bson:"amount,omitempty"`
	Proofs       *[]string `json:"proofs,omitempty" bson:"proofs,omitempty"`

	OR []bson.M `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type MerkledropProofCreate struct {
	ID           *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MerkledropID int64               `json:"merkledrop_id" bson:"merkledrop_id,omitempty"`
	Address      string              `json:"address" bson:"address,omitempty"`
	Index        int64               `json:"index" bson:"index"`
	Amount       int64               `json:"amount" bson:"amount,omitempty"`
	Proofs       []string            `json:"proofs" bson:"proofs,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty" validate:"required"`
}

/**
 * OPERATIONS
 */

// Read

func (m *MerkledropProof) One(filter *MerkledropProofWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *MerkledropProof) List(filter *MerkledropProofWhere, orderBy *MerkledropProofOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*MerkledropProof, error) {
	var items []*MerkledropProof
	orderByKey := "created_at"
	orderByValue := -1

	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
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

func (m *MerkledropProof) Count(filter *MerkledropProofWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *MerkledropProof) Create(data *MerkledropProofCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: checkPools unique
	item := new(MerkledropProof)
	f := bson.M{
		"$and": []bson.M{
			{"merkledrop_id": data.MerkledropID},
			{"address": data.Address},
			{"index": data.Index},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)

	if !item.ID.IsZero() {
		return errors.New("item already stored")
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
