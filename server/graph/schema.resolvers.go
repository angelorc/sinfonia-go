package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	simodel "github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Transaction(ctx context.Context, where *simodel.TransactionWhere) (*simodel.Transaction, error) {
	if where == nil {
		where = &simodel.TransactionWhere{}
	}

	item := simodel.Transaction{}
	item.One(where)
	if item.Hash == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Transactions(ctx context.Context, where *simodel.TransactionWhere, in []*primitive.ObjectID, orderBy *simodel.TransactionOrderByENUM, skip *int, limit *int) ([]*simodel.Transaction, error) {
	if where == nil {
		where = &simodel.TransactionWhere{}
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := simodel.Transaction{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) TransactionCount(ctx context.Context, where *simodel.TransactionWhere) (*int, error) {
	t := simodel.Transaction{}
	if where == nil {
		where = &simodel.TransactionWhere{}
	}

	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) Message(ctx context.Context, where *simodel.MessageWhere) (*simodel.Message, error) {
	if where == nil {
		where = &simodel.MessageWhere{}
	}

	item := simodel.Message{}
	item.One(where)
	if item.ChainID == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Messages(ctx context.Context, where *simodel.MessageWhere, in []*primitive.ObjectID, orderBy *simodel.MessageOrderByENUM, skip *int, limit *int) ([]*simodel.Message, error) {
	if where == nil {
		where = &simodel.MessageWhere{}
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := simodel.Message{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) MessageCount(ctx context.Context, where *simodel.MessageWhere) (*int, error) {
	t := simodel.Message{}
	if where == nil {
		where = &simodel.MessageWhere{}
	}

	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) Account(ctx context.Context, where *simodel.AccountWhere) (*simodel.Account, error) {
	if where == nil {
		where = &simodel.AccountWhere{}
	}

	item := simodel.Account{}
	item.One(where)
	if item.Address == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Accounts(ctx context.Context, where *simodel.AccountWhere, in []*primitive.ObjectID, orderBy *simodel.AccountOrderByENUM, skip *int, limit *int) ([]*simodel.Account, error) {
	if where == nil {
		where = &simodel.AccountWhere{}
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := simodel.Account{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) AccountCount(ctx context.Context, where *simodel.AccountWhere) (*int, error) {
	t := simodel.Account{}
	if where == nil {
		where = &simodel.AccountWhere{}
	}

	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) Incentive(ctx context.Context, where *simodel.IncentiveWhere) (*simodel.Incentive, error) {
	if where == nil {
		where = &simodel.IncentiveWhere{}
	}

	item := simodel.Incentive{}
	item.One(where)
	if item.Receiver == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Incentives(ctx context.Context, where *simodel.IncentiveWhere, in []*primitive.ObjectID, orderBy *simodel.IncentiveOrderByENUM, skip *int, limit *int) ([]*simodel.Incentive, error) {
	if where == nil {
		where = &simodel.IncentiveWhere{}
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := simodel.Incentive{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) IncentiveCount(ctx context.Context, where *simodel.IncentiveWhere) (*int, error) {
	t := simodel.Incentive{}
	if where == nil {
		where = &simodel.IncentiveWhere{}
	}

	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) Swap(ctx context.Context, where *simodel.SwapWhere) (*simodel.Swap, error) {
	if where == nil {
		where = &simodel.SwapWhere{}
	}

	item := simodel.Swap{}
	item.One(where)
	if item.TxHash == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Swaps(ctx context.Context, where *simodel.SwapWhere, in []*primitive.ObjectID, orderBy *simodel.SwapOrderByENUM, skip *int, limit *int) ([]*simodel.Swap, error) {
	if where == nil {
		where = &simodel.SwapWhere{}
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := simodel.Swap{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) SwapCount(ctx context.Context, where *simodel.SwapWhere) (*int, error) {
	t := simodel.Swap{}
	if where == nil {
		where = &simodel.SwapWhere{}
	}

	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *queryResolver) Pool(ctx context.Context, where *simodel.PoolWhere) (*simodel.Pool, error) {
	if where == nil {
		where = &simodel.PoolWhere{}
	}

	item := simodel.Pool{}
	item.One(where)
	if item.TxHash == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Pools(ctx context.Context, where *simodel.PoolWhere, in []*primitive.ObjectID, orderBy *simodel.PoolOrderByENUM, skip *int, limit *int) ([]*simodel.Pool, error) {
	if where == nil {
		where = &simodel.PoolWhere{}
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := simodel.Pool{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) PoolCount(ctx context.Context, where *simodel.PoolWhere) (*int, error) {
	t := simodel.Pool{}
	if where == nil {
		where = &simodel.PoolWhere{}
	}

	count, err := t.Count(where)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
