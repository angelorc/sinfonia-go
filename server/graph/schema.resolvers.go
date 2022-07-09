package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/angelorc/sinfonia-go/config"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
	"github.com/angelorc/sinfonia-go/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) UpdateMerkledrop(ctx context.Context, id primitive.ObjectID, data model.MerkledropUpdateReq) (*model.Merkledrop, error) {
	item := model.Merkledrop{}

	if utility.IsZeroVal(id) {
		return &model.Merkledrop{}, errors.New("missing merkledrop id")
	}

	// Validate
	if err := utility.ValidateStruct(data); err != nil {
		return &model.Merkledrop{}, err
	}

	dataUpdate := model.MerkledropUpdate{}
	dataUpdate.Name = data.Name

	// Upload
	if data.Image != nil {
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		body := &bytes.Buffer{}
		bodywriter := multipart.NewWriter(body)

		writer, err := bodywriter.CreateFormFile("file", data.Name)
		if err != nil {
			return &model.Merkledrop{}, err
		}

		_, err = io.Copy(writer, data.Image.File)
		if err != nil {
			return &model.Merkledrop{}, err
		}

		err = bodywriter.Close()
		if err != nil {
			return &model.Merkledrop{}, err
		}

		cloudFlareImagesUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/images/v1", c.GetSecret("CLOUDFLARE_ACCOUNT"))
		req, err := http.NewRequest("POST", cloudFlareImagesUrl, bytes.NewReader(body.Bytes()))
		if err != nil {
			return &model.Merkledrop{}, err
		}

		req.Header.Set("Content-Type", bodywriter.FormDataContentType())
		req.Header.Add("Authorization", "Bearer "+c.GetSecret("CLOUDFLARE_IMAGES"))
		rsp, _ := client.Do(req)
		if rsp.StatusCode != http.StatusOK {
			return &model.Merkledrop{}, fmt.Errorf("request failed with response code: %d", rsp.StatusCode)
		}
		defer rsp.Body.Close()
		rspBz, _ := ioutil.ReadAll(rsp.Body)

		var cloudlfareResp model.MerkledropUpdateImageResponse
		if err := json.Unmarshal(rspBz, &cloudlfareResp); err != nil {
			return &model.Merkledrop{}, err
		}

		if cloudlfareResp.Success {
			if len(cloudlfareResp.Result.Variants) > 0 {
				dataUpdate.Image = &cloudlfareResp.Result.Variants[0]
			}
		}
	}

	if err := item.Update(id, &dataUpdate); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *mutationResolver) StoreMerkledropProofs(ctx context.Context, id int, file graphql.Upload) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid merkledrop_id")
	}
	if file.Size <= 0 {
		return 0, nil
	}

	parsedList, err := parseMerkleProofsList(file.File)
	if err != nil {
		return 0, err
	}

	total := 0

	// TODO: store proofs in batch mode
	for addr, r := range parsedList {
		item := model.MerkledropProof{}

		amount, err := strconv.ParseInt(r.Amount, 10, 64)
		if err != nil {
			return 0, err
		}

		data := model.MerkledropProofCreate{
			MerkledropID: int64(id),
			Address:      addr,
			Index:        r.Index,
			Amount:       amount,
			Proofs:       r.Proof,
			CreatedAt:    time.Now(),
		}

		if err := item.Create(&data); err != nil {
			return total, err
		}

		total += 1
	}

	return total, nil
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
