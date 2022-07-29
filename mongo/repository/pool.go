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
	poolCollectionName = "pools"
	poolDbRefName      = "default"
)

type poolRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type PoolRepository interface {
	Count(filter *modelv2.PoolFilter) (int64, error)
	Find(filter *modelv2.PoolFilter, pagination *types.PaginationReq) ([]*modelv2.Pool, error)
	FindOne(filter *modelv2.PoolFilter) *modelv2.Pool
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.Pool
	FindByPoolID(poolID uint64) *modelv2.Pool

	Create(data *modelv2.PoolCreateReq) (*primitive.ObjectID, error)
}

func NewPoolRepository() PoolRepository {
	coll := db.GetCollection(poolCollectionName, poolDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	//defer cancel()

	return &poolRepository{context: ctx, collection: coll}
}

func (e *poolRepository) FindOne(filter *modelv2.PoolFilter) *modelv2.Pool {
	var pool modelv2.Pool
	e.collection.FindOne(e.context, &filter).Decode(&pool)

	return &pool
}

func (e *poolRepository) FindByID(id primitive.ObjectID) *modelv2.Pool {
	return e.FindOne(&modelv2.PoolFilter{Id: &id})
}

func (e *poolRepository) FindByPoolID(poolID uint64) *modelv2.Pool {
	return e.FindOne(&modelv2.PoolFilter{PoolID: &poolID})
}

func (e *poolRepository) Find(filter *modelv2.PoolFilter, pagination *types.PaginationReq) ([]*modelv2.Pool, error) {
	var pools []*modelv2.Pool

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
		return pools, err
	}
	err = cursor.All(e.context, &pools)
	if err != nil {
		return pools, err
	}

	return pools, nil
}

func (e *poolRepository) Count(filter *modelv2.PoolFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *poolRepository) Create(data *modelv2.PoolCreateReq) (*primitive.ObjectID, error) {
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

func (e *poolRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"height", -1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"pool_id", -1},
		},
		Options: options.Index().SetUnique(true),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}
