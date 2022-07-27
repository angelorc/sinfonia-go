package repository

import (
	"context"
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
	fantokenCollectionName = "fantokens"
	fantokenDbRefName      = "default"
)

type fantokenRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type FantokenRepository interface {
	Count(filter *types.FantokenFilter) (int64, error)
	Find(filter *types.FantokenFilter, pagination *types.PaginationReq) ([]*modelv2.Fantoken, error)
	FindOne(filter *types.FantokenFilter) *modelv2.Fantoken
	EnsureIndexes() ([]string, error)

	FindByID(id primitive.ObjectID) *modelv2.Fantoken
	FindByHeight(height int64) *modelv2.Fantoken
	FindByDenom(denom string) *modelv2.Fantoken

	Create(data *types.FantokenCreateReq) (*modelv2.Fantoken, error)

	Earliest() *modelv2.Fantoken
	Latest() *modelv2.Fantoken
}

func NewFantokenRepository() FantokenRepository {
	coll := db.GetCollection(blockCollectionName, blockDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	return &fantokenRepository{context: ctx, collection: coll}
}

func (f fantokenRepository) Count(filter *types.FantokenFilter) (int64, error) {
	return f.collection.CountDocuments(f.context, &filter)
}

func (f fantokenRepository) Find(filter *types.FantokenFilter, pagination *types.PaginationReq) ([]*modelv2.Fantoken, error) {
	var fantokens []*modelv2.Fantoken

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

	cursor, err := f.collection.Find(f.context, &queryFilter, options)
	if err != nil {
		return fantokens, err
	}
	err = cursor.All(f.context, &fantokens)
	if err != nil {
		return fantokens, err
	}

	return fantokens, nil
}

func (f fantokenRepository) FindOne(filter *types.FantokenFilter) *modelv2.Fantoken {
	var fantoken modelv2.Fantoken
	f.collection.FindOne(f.context, &filter).Decode(&fantoken)

	return &fantoken
}

func (f fantokenRepository) EnsureIndexes() ([]string, error) {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"height", -1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{"denom", 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	return f.collection.Indexes().CreateMany(f.context, indexes)
}

func (f fantokenRepository) FindByID(id primitive.ObjectID) *modelv2.Fantoken {
	//TODO implement me
	panic("implement me")
}

func (f fantokenRepository) FindByHeight(height int64) *modelv2.Fantoken {
	//TODO implement me
	panic("implement me")
}

func (f fantokenRepository) FindByDenom(denom string) *modelv2.Fantoken {
	//TODO implement me
	panic("implement me")
}

func (f fantokenRepository) Create(data *types.FantokenCreateReq) (*modelv2.Fantoken, error) {
	//TODO implement me
	panic("implement me")
}

func (f fantokenRepository) Earliest() *modelv2.Fantoken {
	//TODO implement me
	panic("implement me")
}

func (f fantokenRepository) Latest() *modelv2.Fantoken {
	//TODO implement me
	panic("implement me")
}
