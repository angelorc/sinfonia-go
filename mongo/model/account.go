package model

import (
	"context"
	"errors"
	"fmt"
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

const DB_COLLECTION_NAME__ACCOUNT = "accounts"
const DB_REF_NAME__ACCOUNT = "default"

/**
 * MODEL
 */

type Account struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Address      string             `json:"address" bson:"address"`
	ValueSwapped string             `json:"value_swapped" bson:"value_swapped"`
	FeesPaid     string             `json:"fees_paid" bson:"fees_paid"`
	TotalTxs     string             `json:"total_txs" bson:"total_txs"`
	FirstSeen    time.Time          `json:"first_seen,omitempty" bson:"first_seen,omitempty" validate:"required"`
}

/**
 * ENUM
 */

type AccountOrderByENUM string

/**
 * DTO
 */

// Read

type AccountWhereUnique struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type AccountWhere struct {
	ID      *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Address *string             `json:"address,omitempty" bson:"address,omitempty"`
	OR      []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write

type AccountCreate struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Address   string              `json:"address" bson:"address"`
	FirstSeen time.Time           `json:"first_seen" bson:"first_seen"`
}

/**
 * OPERATIONS
 */

// Read

func (m *Account) One(filter *AccountWhere) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__ACCOUNT, DB_REF_NAME__ACCOUNT)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)

	return nil
}

func (m *Account) List(filter *AccountWhere, orderBy *AccountOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*Account, error) {
	var items []*Account
	orderByKey := "first_seen"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__ACCOUNT, DB_REF_NAME__ACCOUNT)
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

func (m *Account) Count(filter *AccountWhere) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__ACCOUNT, DB_REF_NAME__ACCOUNT)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

func (m *Account) Create(data *AccountCreate) error {
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__ACCOUNT, DB_REF_NAME__ACCOUNT)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: check unique
	item := new(Account)
	f := bson.M{
		"$and": []bson.M{
			{"address": data.Address},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Address != "" {
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

func EnsureAccount(acc string, firstSeen time.Time) error {
	item := Account{}
	data := AccountCreate{
		Address:   acc,
		FirstSeen: firstSeen,
	}

	if err := utility.ValidateStruct(data); err != nil {
		return err
	}

	if err := item.Create(&data); err != nil {
		return err
	}

	return nil
}

func SyncAccounts() error {
	lasBlock := GetLastHeight("osmosis-1")

	// get last block synced from account
	sync := new(Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Accounts = int64(0)
	}

	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__MESSAGE, DB_REF_NAME__MESSAGE)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"height": bson.M{
					"$gt":  sync.Accounts,
					"$lte": lasBlock,
				},
			},
		},
		{
			"$project": bson.M{
				"signer":     1,
				"first_seen": "$time",
			},
		},
		{
			"$group": bson.M{
				"_id": "$signer",
				"first_seen": bson.M{
					"$first": "$first_seen",
				},
			},
		},
	}

	// record struct
	type acc struct {
		Signer    string    `bson:"_id,omitempty"`
		FirstSeen time.Time `bson:"first_seen,omitempty"`
	}

	// aggregate pipeline
	accCursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	// decode
	var accounts []acc
	if err = accCursor.All(ctx, &accounts); err != nil {
		return err
	}

	for _, a := range accounts {
		if err := EnsureAccount(a.Signer, a.FirstSeen); err != nil {
			return err
		}
	}

	// update sync with last synced height
	sync.Accounts = lasBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%d accounts synced to block %d ", len(accounts), sync.Accounts))

	return nil
}
