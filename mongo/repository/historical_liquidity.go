package repository

import (
	"context"
	"fmt"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	historicalLiquidityCollectionName = "historical_liquidity"
	historicalLiquidityDbRefName      = "default"
)

type historicalLiquidityRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type HistoricalLiquidityRepository interface {
	EnsureIndexes() (string, error)
	Create(data *modelv2.HistoricalLiquidityCreateReq) (*primitive.ObjectID, error)
}

func NewHistoricalLiquidityRepository() HistoricalLiquidityRepository {
	coll := db.GetCollection(historicalLiquidityCollectionName, historicalLiquidityDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Second)
	//defer cancel()

	return &historicalLiquidityRepository{context: ctx, collection: coll}
}

func (e *historicalLiquidityRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"pool_id", 1},
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

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"time", -1},
		},
		Options: options.Index().SetUnique(false),
	}

	e.collection.Indexes().CreateOne(e.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"pool_id", 1},
			{"time", -1},
		},
		Options: options.Index().SetUnique(true),
	}

	return e.collection.Indexes().CreateOne(e.context, index)
}

func (e *historicalLiquidityRepository) Create(data *modelv2.HistoricalLiquidityCreateReq) (*primitive.ObjectID, error) {
	data.ID = primitive.NewObjectID()

	/*if err := data.Validate(); err != nil {
		return &primitive.ObjectID{}, err
	}*/

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
