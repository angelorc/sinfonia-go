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
	liquidityEventCollectionName = "liquidity_events"
	liquidityEventDbRefName      = "default"
)

type liquidityEventRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type LiquidityRepository interface {
	Count(filter *modelv2.LiquidityEventFilter) (int64, error)
	Find(filter *modelv2.LiquidityEventFilter, pagination *types.PaginationReq) ([]*modelv2.LiquidityEvent, error)
	FindOne(filter *modelv2.LiquidityEventFilter) *modelv2.LiquidityEvent
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.LiquidityEvent
	FindByHeight(height int64) []*modelv2.LiquidityEvent
	FindBySender(sender string) []*modelv2.LiquidityEvent

	Create(data *modelv2.LiquidityEventCreateReq) (*primitive.ObjectID, error)
}

func NewLiquidityRepository() LiquidityRepository {
	coll := db.GetCollection(liquidityEventCollectionName, liquidityEventDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)
	//defer cancel()

	return &liquidityEventRepository{context: ctx, collection: coll}
}

func (e *liquidityEventRepository) FindOne(filter *modelv2.LiquidityEventFilter) *modelv2.LiquidityEvent {
	var result modelv2.LiquidityEvent
	e.collection.FindOne(e.context, &filter).Decode(&result)

	return &result
}

func (e *liquidityEventRepository) FindByID(id primitive.ObjectID) *modelv2.LiquidityEvent {
	return e.FindOne(&modelv2.LiquidityEventFilter{Id: &id})
}

func (e *liquidityEventRepository) FindByHeight(height int64) []*modelv2.LiquidityEvent {
	results, _ := e.Find(&modelv2.LiquidityEventFilter{Height: &height}, &types.PaginationReq{})
	return results
}

func (e *liquidityEventRepository) Find(filter *modelv2.LiquidityEventFilter, pagination *types.PaginationReq) ([]*modelv2.LiquidityEvent, error) {
	var results []*modelv2.LiquidityEvent

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
		return results, err
	}
	err = cursor.All(e.context, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (e *liquidityEventRepository) FindBySender(sender string) []*modelv2.LiquidityEvent {
	results, _ := e.Find(&modelv2.LiquidityEventFilter{Sender: &sender}, &types.PaginationReq{})
	return results
}

func (e *liquidityEventRepository) Count(filter *modelv2.LiquidityEventFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *liquidityEventRepository) Create(data *modelv2.LiquidityEventCreateReq) (*primitive.ObjectID, error) {
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

func (e *liquidityEventRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"height", -1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"sender", 1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"height", 1},
			{"tx_hash", 1},
			{"sender", 1},
			{"pool_id", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}
