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
	eventCollectionName = "events"
	eventDbRefName      = "default"
)

type eventRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type EventRepository interface {
	Count(filter *modelv2.EventFilter) (int64, error)
	Find(filter *modelv2.EventFilter, pagination *types.PaginationReq) ([]*modelv2.Event, error)
	FindOne(filter *modelv2.EventFilter) *modelv2.Event
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.Event

	Create(data *modelv2.EventCreateReq) (*modelv2.Event, error)
}

func NewEventRepository() EventRepository {
	coll := db.GetCollection(eventCollectionName, eventDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	return &eventRepository{context: ctx, collection: coll}
}

func (e *eventRepository) FindOne(filter *modelv2.EventFilter) *modelv2.Event {
	var event modelv2.Event
	e.collection.FindOne(e.context, &filter).Decode(&event)

	return &event
}

func (e *eventRepository) FindByID(id primitive.ObjectID) *modelv2.Event {
	return e.FindOne(&modelv2.EventFilter{Id: &id})
}

func (e *eventRepository) Find(filter *modelv2.EventFilter, pagination *types.PaginationReq) ([]*modelv2.Event, error) {
	var events []*modelv2.Event

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
		return events, err
	}
	err = cursor.All(e.context, &events)
	if err != nil {
		return events, err
	}

	return events, nil
}

func (e *eventRepository) Count(filter *modelv2.EventFilter) (int64, error) {
	return e.collection.CountDocuments(e.context, &filter)
}

func (e *eventRepository) Create(data *modelv2.EventCreateReq) (*modelv2.Event, error) {
	data.ID = primitive.NewObjectID()

	if err := data.Validate(); err != nil {
		return &modelv2.Event{}, err
	}

	res, err := e.collection.InsertOne(e.context, &data)
	if err != nil {
		return &modelv2.Event{}, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return &modelv2.Event{}, fmt.Errorf("server error")
	}

	return e.FindByID(insertedID), nil
}

func (e *eventRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"tx_id", 1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"tx_id", 1},
			{"msg_index", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"type", 1},
		},
		Options: options.Index().SetUnique(false),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}
