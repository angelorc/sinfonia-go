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
	incentiveCollectionName = "incentives"
	incentiveDbRefName      = "default"
)

type incentiveRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type IncentiveRepository interface {
	Count(filter *modelv2.IncentiveFilter) (int64, error)
	Find(filter *modelv2.IncentiveFilter, pagination *types.PaginationReq) ([]*modelv2.Incentive, error)
	FindOne(filter *modelv2.IncentiveFilter) *modelv2.Incentive
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.Incentive
	FindByHeight(height int64) *modelv2.Incentive
	FindByReceiver(receiver string) []*modelv2.Incentive

	Create(data *modelv2.IncentiveCreateReq) (*primitive.ObjectID, error)
	CreateMany(data []*modelv2.IncentiveCreateReq) (bool, error)
}

func NewIncentiveRepository() IncentiveRepository {
	coll := db.GetCollection(incentiveCollectionName, incentiveDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)
	//defer cancel()

	return &incentiveRepository{context: ctx, collection: coll}
}

func (e *incentiveRepository) FindOne(filter *modelv2.IncentiveFilter) *modelv2.Incentive {
	var incentive modelv2.Incentive
	e.collection.FindOne(e.context, &filter).Decode(&incentive)

	return &incentive
}

func (e *incentiveRepository) FindByID(id primitive.ObjectID) *modelv2.Incentive {
	return e.FindOne(&modelv2.IncentiveFilter{Id: &id})
}

func (e *incentiveRepository) FindByHeight(height int64) *modelv2.Incentive {
	return e.FindOne(&modelv2.IncentiveFilter{Height: &height})
}

func (e *incentiveRepository) Find(filter *modelv2.IncentiveFilter, pagination *types.PaginationReq) ([]*modelv2.Incentive, error) {
	var incentives []*modelv2.Incentive

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
		return incentives, err
	}
	err = cursor.All(e.context, &incentives)
	if err != nil {
		return incentives, err
	}

	return incentives, nil
}

func (e *incentiveRepository) FindByReceiver(receiver string) []*modelv2.Incentive {
	incentives, _ := e.Find(&modelv2.IncentiveFilter{Receiver: &receiver}, &types.PaginationReq{})
	return incentives
}

func (e *incentiveRepository) Count(filter *modelv2.IncentiveFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *incentiveRepository) Create(data *modelv2.IncentiveCreateReq) (*primitive.ObjectID, error) {
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

func (e *incentiveRepository) CreateMany(data []*modelv2.IncentiveCreateReq) (bool, error) {
	newValues := make([]interface{}, len(data))

	for _, record := range data {
		record.ID = primitive.NewObjectID()
		if err := record.Validate(); err != nil {
			return false, err
		}

		newValues = append(newValues, record)
	}

	_, err := e.collection.InsertMany(e.context, newValues)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (e *incentiveRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"height", -1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"receiver", 1},
		},
		Options: options.Index().SetUnique(false),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}
