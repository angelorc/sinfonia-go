package model

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
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
	Claimed      bool     `json:"claimed" bson:"claimed"`

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
	Claimed      bool                `json:"claimed" bson:"claimed"`
	CreatedAt    time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty" validate:"required"`
}

type MerkledropProofClaim struct {
	MerkledropID int64  `json:"merkledrop_id" bson:"merkledrop_id,omitempty" validate:"required"`
	Address      string `json:"address" bson:"address,omitempty" validate:"required"`
	Index        int64  `json:"index" bson:"index" validate:"required"`
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

func (m *MerkledropProof) StoreMany(data []MerkledropProofCreate) error {
	// validate
	for _, el := range data {
		if err := utility.ValidateStruct(el); err != nil {
			return err
		}
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	items := make([]interface{}, len(data))

	for i, proof := range data {
		items[i] = proof
	}

	// operation
	_, err := collection.InsertMany(ctx, items)
	if err != nil {
		if strings.HasPrefix(err.Error(), "bulk write exception") {
			return fmt.Errorf("duplicate record")
		}

		return err
	}

	return nil
}

func (m *MerkledropProof) SetClaimed(data MerkledropProofClaim) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// get
	filter := MerkledropProofWhere{
		MerkledropID: &data.MerkledropID,
		Address:      &data.Address,
		Index:        &data.Index,
	}
	collection.FindOne(ctx, &filter).Decode(&m)

	if m.ID.IsZero() {
		return fmt.Errorf("claim not found")
	}

	// update
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"claimed": true}})
	if err != nil {
		return err
	}

	collection.FindOne(ctx, bson.M{"merkledrop_id": data.MerkledropID}).Decode(&m)

	return nil
}

func (m *MerkledropProof) CreateIndexes() error {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"merkledrop_id", 1},
			{"address", 1},
			{"index", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MERKLEDROP_PROOF, DB_REF_NAME__MERKLEDROP_PROOF)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return fmt.Errorf("error while creting indexes on merkledrop_proofs: %v", err)
	}

	return nil
}
