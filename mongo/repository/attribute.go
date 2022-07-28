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
	attributeCollectionName = "attributess"
	attributeDbRefName      = "default"
)

type attributeRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type AttributeRepository interface {
	Count(filter *types.AttributeFilter) (int64, error)
	Find(filter *types.AttributeFilter, pagination *types.PaginationReq) ([]*modelv2.Attribute, error)
	FindOne(filter *types.AttributeFilter) *modelv2.Attribute
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.Attribute
	FindByKey(key string) *modelv2.Attribute

	Create(data *types.AttributeCreateReq) (*modelv2.Attribute, error)
}

func NewAttributeRepository() AttributeRepository {
	coll := db.GetCollection(attributeCollectionName, attributeDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	return &attributeRepository{context: ctx, collection: coll}
}

func (e *attributeRepository) FindOne(filter *types.AttributeFilter) *modelv2.Attribute {
	var attribute modelv2.Attribute
	e.collection.FindOne(e.context, &filter).Decode(&attribute)

	return &attribute
}

func (e *attributeRepository) FindByID(id primitive.ObjectID) *modelv2.Attribute {
	return e.FindOne(&types.AttributeFilter{Id: &id})
}

func (e *attributeRepository) FindByKey(key string) *modelv2.Attribute {
	return e.FindOne(&types.AttributeFilter{Key: &key})
}

func (e *attributeRepository) Find(filter *types.AttributeFilter, pagination *types.PaginationReq) ([]*modelv2.Attribute, error) {
	var attributes []*modelv2.Attribute

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
		return attributes, err
	}
	err = cursor.All(e.context, &attributes)
	if err != nil {
		return attributes, err
	}

	return attributes, nil
}

func (e *attributeRepository) Count(filter *types.AttributeFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *attributeRepository) Create(data *types.AttributeCreateReq) (*modelv2.Attribute, error) {
	data.ID = primitive.NewObjectID()

	if err := data.Validate(); err != nil {
		return &modelv2.Attribute{}, err
	}

	res, err := e.collection.InsertOne(e.context, &data)
	if err != nil {
		return &modelv2.Attribute{}, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return &modelv2.Attribute{}, fmt.Errorf("server error")
	}

	return e.FindByID(insertedID), nil
}

func (e *attributeRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"event_id", 1},
			{"key", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}
