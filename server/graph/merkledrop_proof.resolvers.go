package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
)

func (r *merkledropProofWhereResolver) Claimed(ctx context.Context, obj *model.MerkledropProofWhere, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// MerkledropProofWhere returns generated.MerkledropProofWhereResolver implementation.
func (r *Resolver) MerkledropProofWhere() generated.MerkledropProofWhereResolver {
	return &merkledropProofWhereResolver{r}
}

type merkledropProofWhereResolver struct{ *Resolver }
