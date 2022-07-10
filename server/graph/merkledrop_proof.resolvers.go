package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
)

func (r *merkledropProofResolver) Merkledrop(ctx context.Context, obj *model.MerkledropProof) (*model.Merkledrop, error) {
	where := &model.MerkledropWhere{
		MerkledropID: &obj.MerkledropID,
	}

	item := model.Merkledrop{}
	item.One(where)
	if item.ID.IsZero() {
		return &model.Merkledrop{}, nil
	}

	return &item, nil
}

func (r *merkledropProofWhereResolver) Claimed(ctx context.Context, obj *model.MerkledropProofWhere, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// MerkledropProof returns generated.MerkledropProofResolver implementation.
func (r *Resolver) MerkledropProof() generated.MerkledropProofResolver {
	return &merkledropProofResolver{r}
}

// MerkledropProofWhere returns generated.MerkledropProofWhereResolver implementation.
func (r *Resolver) MerkledropProofWhere() generated.MerkledropProofWhereResolver {
	return &merkledropProofWhereResolver{r}
}

type merkledropProofResolver struct{ *Resolver }
type merkledropProofWhereResolver struct{ *Resolver }
