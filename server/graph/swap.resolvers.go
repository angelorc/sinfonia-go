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

func (r *swapResolver) TxHash(ctx context.Context, obj *model.Swap) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *swapResolver) Timestamp(ctx context.Context, obj *model.Swap) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *swapWhereResolver) TxHash(ctx context.Context, obj *model.SwapWhere, data *string) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *swapWhereResolver) Timestamp(ctx context.Context, obj *model.SwapWhere, data *time.Time) error {
	panic(fmt.Errorf("not implemented"))
}

// Swap returns generated.SwapResolver implementation.
func (r *Resolver) Swap() generated.SwapResolver { return &swapResolver{r} }

// SwapWhere returns generated.SwapWhereResolver implementation.
func (r *Resolver) SwapWhere() generated.SwapWhereResolver { return &swapWhereResolver{r} }

type swapResolver struct{ *Resolver }
type swapWhereResolver struct{ *Resolver }
