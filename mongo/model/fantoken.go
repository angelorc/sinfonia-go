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

func SyncFantokens() error {
	lastBlock := GetLastHeight()

	// get last block synced from account
	sync := new(Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Fantokens = int64(0)
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
					"$gt":  sync.Fantokens,
					"$lte": lastBlock,
				},
			},
		},
		{
			"$match": bson.M{
				"msg_type": bson.M{
					"$eq": "/bitsong.fantoken.MsgIssueFanToken",
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
			"$project": bson.M{
				"chain_id":      1,
				"height":        1,
				"tx_id":         1,
				"signer":        1,
				"time":          1,
				"tx.log.events": 1,
			},
		},
	}

	// record struct
	type attribute struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
	type event struct {
		Type       string      `bson:"type"`
		Attributes []attribute `bson:"attributes"`
	}
	type txlog struct {
		Events []event `bson:"events"`
	}
	type tx struct {
		Log []txlog `bson:"log"`
	}
	type record struct {
		Signer  string             `bson:"signer"`
		Time    time.Time          `bson:"time"`
		ChainID string             `json:"chain_id" bson:"chain_id" validate:"required"`
		Height  int64              `json:"height" bson:"height" validate:"required"`
		TxID    primitive.ObjectID `json:"tx_id" bson:"tx_id" validate:"required"`
		Tx      []tx               `bson:"tx"`
	}

	// aggregate pipeline
	accCursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	// decode
	var records []record
	if err = accCursor.All(ctx, &records); err != nil {
		return err
	}

	for _, r := range records {
		for _, rtx := range r.Tx {
			for _, rtxlog := range rtx.Log {
				for _, evt := range rtxlog.Events {
					denom := ""

					switch evt.Type {
					case "issue_fantoken":
						denom = evt.Attributes[0].Value
					}

					if denom != "" {
						fantoken := new(Fantoken)
						data := &FantokenCreate{
							ChainID:  &r.ChainID,
							Height:   &r.Height,
							TxID:     &r.TxID,
							Denom:    &denom,
							Owner:    &r.Signer,
							IssuedAt: &r.Time,
						}

						if err := fantoken.Create(data); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// update sync with last synced height
	sync.Fantokens = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%d fantokens synced to block %d ", len(records), sync.Fantokens))

	return nil
}
