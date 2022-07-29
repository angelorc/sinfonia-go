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
	transactionCollectionName = "transaction"
	transactionDbRefName      = "default"
)

type transactionRepository struct {
	context    context.Context
	collection *mongo.Collection
}

type TransactionRepository interface {
	Count(filter *types.TransactionFilter) (int64, error)
	Find(filter *types.TransactionFilter, pagination *types.PaginationReq) ([]*modelv2.Transaction, error)
	FindOne(filter *types.TransactionFilter) *modelv2.Transaction
	EnsureIndexes() (string, error)

	FindByID(id primitive.ObjectID) *modelv2.Transaction
	FindByHash(hash string) *modelv2.Transaction

	Create(data *types.TransactionCreateReq) (*modelv2.Transaction, error)
}

func NewTransactionRepository() TransactionRepository {
	coll := db.GetCollection(transactionCollectionName, transactionDbRefName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	return &transactionRepository{context: ctx, collection: coll}
}

func (b *transactionRepository) FindOne(filter *types.TransactionFilter) *modelv2.Transaction {
	var transaction modelv2.Transaction
	b.collection.FindOne(b.context, &filter).Decode(&transaction)

	return &transaction
}

func (b *transactionRepository) FindByID(id primitive.ObjectID) *modelv2.Transaction {
	return b.FindOne(&types.TransactionFilter{Id: &id})
}

func (b *transactionRepository) FindByHash(hash string) *modelv2.Transaction {
	return b.FindOne(&types.TransactionFilter{Hash: &hash})
}

func (b *transactionRepository) Find(filter *types.TransactionFilter, pagination *types.PaginationReq) ([]*modelv2.Transaction, error) {
	var transactions []*modelv2.Transaction

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
		return transactions, err
	}
	err = cursor.All(b.context, &transactions)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (b *transactionRepository) Count(filter *types.TransactionFilter) (int64, error) {
	return b.collection.CountDocuments(b.context, &filter)
}

func (b *transactionRepository) Create(data *types.TransactionCreateReq) (*modelv2.Transaction, error) {
	// create id
	txID, err := primitive.ObjectIDFromHex(data.Hash[:24])
	if err != nil {
		panic(err)
	}
	data.ID = txID

	if err := data.Validate(); err != nil {
		return &modelv2.Transaction{}, err
	}

	res, err := b.collection.InsertOne(b.context, &data)
	if err != nil {
		return &modelv2.Transaction{}, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return &modelv2.Transaction{}, fmt.Errorf("server error")
	}
	fmt.Println(insertedID)
	return b.FindByID(insertedID), nil
}

func (b *transactionRepository) EnsureIndexes() (string, error) {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"height", -1},
		},
		Options: options.Index().SetUnique(false),
	}

	b.collection.Indexes().CreateOne(b.context, index)

	index = mongo.IndexModel{
		Keys: bson.D{
			{"hash", 1},
		},
		Options: options.Index().SetUnique(true),
	}

	return b.collection.Indexes().CreateOne(b.context, index)
}
