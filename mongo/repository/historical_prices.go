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
	historicalPricesCollectionName = "historical_prices"
	historicalPricesDbRefName      = "default"
)

type historicalPriceRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type HistoricalPriceRepository interface {
	Count(filter *modelv2.HistoricalPriceFilter) (int64, error)
	Find(filter *modelv2.HistoricalPriceFilter, pagination *types.PaginationReq) ([]*modelv2.HistoricalPrice, error)
	FindOne(filter *modelv2.HistoricalPriceFilter) *modelv2.HistoricalPrice
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.HistoricalPrice
	FindByAsset(asset string) []*modelv2.HistoricalPrice

	Create(data *modelv2.HistoricalPriceCreateReq) (*primitive.ObjectID, error)
}

func NewHistoricalPriceRepository() HistoricalPriceRepository {
	coll := db.GetCollection(historicalPricesCollectionName, historicalPricesDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 1200*time.Second)
	//defer cancel()

	return &historicalPriceRepository{context: ctx, collection: coll}
}

func (e *historicalPriceRepository) Count(filter *modelv2.HistoricalPriceFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *historicalPriceRepository) FindOne(filter *modelv2.HistoricalPriceFilter) *modelv2.HistoricalPrice {
	var hp modelv2.HistoricalPrice
	e.collection.FindOne(e.context, &filter).Decode(&hp)

	return &hp
}

func (e *historicalPriceRepository) Find(filter *modelv2.HistoricalPriceFilter, pagination *types.PaginationReq) ([]*modelv2.HistoricalPrice, error) {
	var hps []*modelv2.HistoricalPrice

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
		return hps, err
	}
	err = cursor.All(e.context, &hps)
	if err != nil {
		return hps, err
	}

	return hps, nil
}

func (e *historicalPriceRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"asset", 1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"time", 1},
		},
		Options: options.Index().SetUnique(false),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}

func (e *historicalPriceRepository) FindByID(id primitive.ObjectID) *modelv2.HistoricalPrice {
	return e.FindOne(&modelv2.HistoricalPriceFilter{Id: &id})
}

func (e *historicalPriceRepository) FindByAsset(asset string) []*modelv2.HistoricalPrice {
	hps, _ := e.Find(
		&modelv2.HistoricalPriceFilter{
			Asset: &asset,
		},
		&types.PaginationReq{},
	)
	return hps
}

func (e *historicalPriceRepository) Create(data *modelv2.HistoricalPriceCreateReq) (*primitive.ObjectID, error) {
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