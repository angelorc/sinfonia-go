package repository

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"github.com/angelorc/sinfonia-go/mongo/types"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	swapCollectionName = "swaps"
	swapDbRefName      = "default"
)

type swapRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type SwapRepository interface {
	Count(filter *modelv2.SwapFilter) (int64, error)
	Find(filter *modelv2.SwapFilter, pagination *types.PaginationReq) ([]*modelv2.Swap, error)
	FindOne(filter *modelv2.SwapFilter) *modelv2.Swap
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.Swap
	FindByHeight(height int64) *modelv2.Swap

	Create(data *modelv2.SwapCreateReq) (*primitive.ObjectID, error)
	InsertMany(records []interface{}) (*mongo.InsertManyResult, error)
}

func NewSwapRepository() SwapRepository {
	// TODO: improve the logic on create and ensure index
	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Second)
	//defer cancel()

	tsOpts := options.TimeSeries().SetTimeField("time")
	tsOpts.SetGranularity("seconds")
	opts := options.CreateCollection().SetTimeSeriesOptions(tsOpts)
	db.GetDB(swapDbRefName).CreateCollection(ctx, swapCollectionName, opts)

	coll := db.GetCollection(swapCollectionName, swapDbRefName)
	swapRepo := &swapRepository{context: ctx, collection: coll}

	swapRepo.EnsureIndexes()

	return swapRepo
}

func (e *swapRepository) FindOne(filter *modelv2.SwapFilter) *modelv2.Swap {
	var swap modelv2.Swap
	e.collection.FindOne(e.context, &filter).Decode(&swap)

	return &swap
}

func (e *swapRepository) FindByID(id primitive.ObjectID) *modelv2.Swap {
	return e.FindOne(&modelv2.SwapFilter{Id: &id})
}

func (e *swapRepository) FindByHeight(height int64) *modelv2.Swap {
	return e.FindOne(&modelv2.SwapFilter{Height: &height})
}

func (e *swapRepository) Find(filter *modelv2.SwapFilter, pagination *types.PaginationReq) ([]*modelv2.Swap, error) {
	var swaps []*modelv2.Swap

	orderByKey := "height"
	orderByValue := -1

	options := options.Find()
	if pagination.Limit != nil {
		options.SetLimit(*pagination.Limit)
	}
	if pagination.Skip != nil {
		options.SetSkip(*pagination.Skip)
	}
	if pagination.OrderBy != nil {
		orderByKey, orderByValue = utility.GetOrderByKeyAndValue(*pagination.OrderBy)
	}
	options.SetSort(map[string]int{orderByKey: orderByValue})

	var queryFilter interface{}
	if filter != nil {
		queryFilter = filter
	}

	cursor, err := e.collection.Find(e.context, &queryFilter, options)
	if err != nil {
		return swaps, err
	}
	err = cursor.All(e.context, &swaps)
	if err != nil {
		return swaps, err
	}

	return swaps, nil
}

func (e *swapRepository) Count(filter *modelv2.SwapFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *swapRepository) Create(data *modelv2.SwapCreateReq) (*primitive.ObjectID, error) {
	data.ID = primitive.NewObjectID()

	if err := data.Validate(); err != nil {
		return &primitive.ObjectID{}, err
	}

	res, err := e.collection.InsertOne(e.context, &data)
	if err != nil {
		return &primitive.ObjectID{}, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return &primitive.ObjectID{}, fmt.Errorf("server error")
	}

	return &insertedID, nil
}

func (e *swapRepository) InsertMany(records []interface{}) (*mongo.InsertManyResult, error) {
	return e.collection.InsertMany(e.context, records)
}

func (e *swapRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"height", -1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"tx_hash", 1},
			{"pool_id", 1},
			{"account", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}
