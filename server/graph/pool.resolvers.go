package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
)

func (r *poolResolver) TxHash(ctx context.Context, obj *model.Pool) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *poolResolver) Timestamp(ctx context.Context, obj *model.Pool) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *poolWhereResolver) TxHash(ctx context.Context, obj *model.PoolWhere, data *string) error {
	panic(fmt.Errorf("not implemented"))
}

// Pool returns generated.PoolResolver implementation.
func (r *Resolver) Pool() generated.PoolResolver { return &poolResolver{r} }

// PoolWhere returns generated.PoolWhereResolver implementation.
func (r *Resolver) PoolWhere() generated.PoolWhereResolver { return &poolWhereResolver{r} }

type poolResolver struct{ *Resolver }
type poolWhereResolver struct{ *Resolver }
