package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

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

func (r *queryResolver) TransactionCount(ctx context.Context, where *simodel.TransactionWhere, search *string) (*int, error) {
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
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Incentives(ctx context.Context, where *simodel.IncentiveWhere, in []*primitive.ObjectID, orderBy *simodel.IncentiveOrderByENUM, skip *int, limit *int) ([]*simodel.Incentive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) IncentiveCount(ctx context.Context, where *simodel.IncentiveWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
