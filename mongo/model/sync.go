package model

import (
	"context"
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

const DB_COLLECTION_NAME__SYNC = "sync"
const DB_REF_NAME__SYNC = "default"

/**
 * MODEL
 */

type Sync struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Accounts  int64              `json:"accounts" bson:"accounts"`
	Fantokens int64              `json:"fantokens" bson:"fantokens"`
}

/**
 * OPERATIONS
 */

// Read

func (s *Sync) One() error {
	collection := db.GetCollection(DB_COLLECTION_NAME__SYNC, DB_REF_NAME__SYNC)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, bson.M{}).Decode(&s)

	return nil
}

// Write Operations

func (s *Sync) Save() error {
	// validate
	if err := utility.ValidateStruct(s); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__SYNC, DB_REF_NAME__SYNC)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// operation
	filter := bson.M{"_id": s.ID}
	update := bson.M{"$set": s}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}

	if err := collection.FindOne(ctx, bson.M{"_id": s.ID}).Decode(&s); err != nil {
		return err
	}

	return nil
}
