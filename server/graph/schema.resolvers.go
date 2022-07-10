package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"github.com/angelorc/sinfonia-go/server/util"
	"strconv"
	"time"

	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) UpdateMerkledrop(ctx context.Context, id int, data model.MerkledropUpdateReq) (*model.Merkledrop, error) {
	item := model.Merkledrop{}

	// Validate
	if id <= 0 {
		return &item, errors.New("invalid merkledrop_id")
	}
	if err := utility.ValidateStruct(data); err != nil {
		return &item, err
	}

	dataUpdate := model.MerkledropUpdate{}
	dataUpdate.Name = data.Name

	// Upload Image
	if data.Image != nil {
		imageUrl, err := util.UploadImage(data)
		if err != nil {
			return &item, err
		}

		dataUpdate.Image = imageUrl
	}

	// Store List
	if data.List != nil && data.List.Size > 0 {
		parsedList, err := parseMerkleProofsList(data.List.File)
		if err != nil {
			return &item, err
		}

		proof := model.MerkledropProof{}
		dataProofs := make([]model.MerkledropProofCreate, 0)

		if err := proof.CreateIndexes(); err != nil {
			return &item, err
		}

		for addr, r := range parsedList {
			amount, err := strconv.ParseInt(r.Amount, 10, 64)
			if err != nil {
				return &item, err
			}

			dataProof := model.MerkledropProofCreate{
				MerkledropID: int64(id),
				Address:      addr,
				Index:        r.Index,
				Amount:       amount,
				Proofs:       r.Proof,
				Claimed:      false,
				CreatedAt:    time.Now(),
			}

			dataProofs = append(dataProofs, dataProof)
		}

		if err := proof.StoreMany(dataProofs); err != nil {
			return &item, err
		}
	}

	if err := item.Update(int64(id), &dataUpdate); err != nil {
		return &item, err
	}

	return &item, nil
}

func (r *queryResolver) Transaction(ctx context.Context, where *model.TransactionWhere) (*model.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Transactions(ctx context.Context, where *model.TransactionWhere, in []*primitive.ObjectID, orderBy *model.TransactionOrderByENUM, skip *int, limit *int) ([]*model.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TransactionCount(ctx context.Context, where *model.TransactionWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Message(ctx context.Context, where *model.MessageWhere) (*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Messages(ctx context.Context, where *model.MessageWhere, in []*primitive.ObjectID, orderBy *model.MessageOrderByENUM, skip *int, limit *int) ([]*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MessageCount(ctx context.Context, where *model.MessageWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Account(ctx context.Context, where *model.AccountWhere) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Accounts(ctx context.Context, where *model.AccountWhere, in []*primitive.ObjectID, orderBy *model.AccountOrderByENUM, skip *int, limit *int) ([]*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AccountCount(ctx context.Context, where *model.AccountWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Fantoken(ctx context.Context, where *model.FantokenWhere) (*model.Fantoken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Fantokens(ctx context.Context, where *model.FantokenWhere, in []*primitive.ObjectID, orderBy *model.FantokenOrderByENUM, skip *int, limit *int) ([]*model.Fantoken, error) {
	if where == nil {
		where = &model.FantokenWhere{}
	}

	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := model.Fantoken{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *queryResolver) FantokenCount(ctx context.Context, where *model.FantokenWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Merkledrop(ctx context.Context, where *model.MerkledropWhere) (*model.Merkledrop, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Merkledrops(ctx context.Context, where *model.MerkledropWhere, in []*primitive.ObjectID, orderBy *model.MerkledropOrderByENUM, skip *int, limit *int) ([]*model.Merkledrop, error) {
	if where == nil {
		where = &model.MerkledropWhere{}
	}

	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := model.Merkledrop{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}

	// check on-chain then remove if is claimed

	return items, nil
}

func (r *queryResolver) MerkledropCount(ctx context.Context, where *model.MerkledropWhere) (*int, error) {
	m := model.Merkledrop{}
	if where == nil {
		where = &model.MerkledropWhere{}
	}
	count, err := m.Count(where)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *queryResolver) MerkledropProof(ctx context.Context, where *model.MerkledropProofWhere) (*model.MerkledropProof, error) {
	if where == nil {
		where = &model.MerkledropProofWhere{}
	}

	item := model.MerkledropProof{}
	item.One(where)
	if item.Address == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) MerkledropProofs(ctx context.Context, where *model.MerkledropProofWhere, in []*primitive.ObjectID, orderBy *model.MerkledropProofOrderByENUM, skip *int, limit *int) ([]*model.MerkledropProof, error) {
	if where == nil {
		where = &model.MerkledropProofWhere{}
	}

	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := model.MerkledropProof{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *queryResolver) MerkledropProofCount(ctx context.Context, where *model.MerkledropProofWhere) (*int, error) {
	m := model.MerkledropProof{}
	if where == nil {
		where = &model.MerkledropProofWhere{}
	}
	count, err := m.Count(where)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *queryResolver) Incentive(ctx context.Context, where *model.IncentiveWhere) (*model.Incentive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Incentives(ctx context.Context, where *model.IncentiveWhere, in []*primitive.ObjectID, orderBy *model.IncentiveOrderByENUM, skip *int, limit *int) ([]*model.Incentive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) IncentiveCount(ctx context.Context, where *model.IncentiveWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Swap(ctx context.Context, where *model.SwapWhere) (*model.Swap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Swaps(ctx context.Context, where *model.SwapWhere, in []*primitive.ObjectID, orderBy *model.SwapOrderByENUM, skip *int, limit *int) ([]*model.Swap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SwapCount(ctx context.Context, where *model.SwapWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pool(ctx context.Context, where *model.PoolWhere) (*model.Pool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pools(ctx context.Context, where *model.PoolWhere, in []*primitive.ObjectID, orderBy *model.PoolOrderByENUM, skip *int, limit *int) ([]*model.Pool, error) {
	if where == nil {
		where = &model.PoolWhere{}
	}

	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := model.Pool{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *queryResolver) PoolCount(ctx context.Context, where *model.PoolWhere) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
