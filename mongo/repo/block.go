package repo

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
	blockCollectionName = "blocks"
	blockDbRefName      = "default"
)

type blockRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type BlockRepository interface {
	Find(filter *types.BlockFilter) *modelv2.Block
	FindByID(id primitive.ObjectID) *modelv2.Block
	FindByHeight(height int64) *modelv2.Block

	List(filter *types.BlockFilter, pagination *types.PaginationReq) ([]*modelv2.Block, error)
	Count(filter *types.BlockFilter) (int64, error)

	EnsureIndexes() (string, error)

	Create(data *types.BlockCreateReq) (*modelv2.Block, error)

	Earliest() *modelv2.Block
	Latest() *modelv2.Block
}

func NewBlockRepository() BlockRepository {
	coll := db.GetCollection(blockCollectionName, blockDbRefName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return &blockRepository{context: ctx, collection: coll}
}

func (b *blockRepository) Find(filter *types.BlockFilter) *modelv2.Block {
	var block modelv2.Block
	b.collection.FindOne(b.context, &filter).Decode(&block)

	return &block
}

func (b *blockRepository) FindByID(id primitive.ObjectID) *modelv2.Block {
	var block modelv2.Block

	filter := types.BlockFilter{Id: &id}
	b.collection.FindOne(b.context, &filter).Decode(&block)

	return &block
}

func (b *blockRepository) FindByHeight(height int64) *modelv2.Block {
	var block modelv2.Block

	filter := types.BlockFilter{Height: &height}
	b.collection.FindOne(b.context, &filter).Decode(&block)

	return &block
}

func (b *blockRepository) List(filter *types.BlockFilter, pagination *types.PaginationReq) ([]*modelv2.Block, error) {
	var blocks []*modelv2.Block

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

	cursor, err := b.collection.Find(b.context, &queryFilter, options)
	if err != nil {
		return blocks, err
	}
	err = cursor.All(b.context, &blocks)
	if err != nil {
		return blocks, err
	}

	return blocks, nil
}

func (b *blockRepository) Count(filter *types.BlockFilter) (int64, error) {
	return b.collection.CountDocuments(b.context, &filter)
}

func (b *blockRepository) Create(data *types.BlockCreateReq) (*modelv2.Block, error) {
	if err := data.Validate(); err != nil {
		return &modelv2.Block{}, err
	}

	res, err := b.collection.InsertOne(b.context, data)
	if err != nil {
		return &modelv2.Block{}, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return &modelv2.Block{}, fmt.Errorf("server error")
	}

	return b.FindByID(insertedID), nil
}

func (b *blockRepository) Earliest() *modelv2.Block {
	var block modelv2.Block

	opts := options.FindOne()
	opts.SetSort(map[string]int{"height": 1})

	b.collection.FindOne(b.context, &types.BlockFilter{}, opts).Decode(block)

	return &block
}

func (b *blockRepository) Latest() *modelv2.Block {
	var block modelv2.Block

	opts := options.FindOne()
	opts.SetSort(map[string]int{"height": -1})

	b.collection.FindOne(b.context, &types.BlockFilter{}, opts).Decode(block)

	return &block
}

func (b *blockRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"height", -1},
		},
		Options: options.Index().SetUnique(true),
	}

	return b.collection.Indexes().CreateOne(b.context, index)
}
