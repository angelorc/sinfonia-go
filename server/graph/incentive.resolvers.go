package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *incentiveAssetWhereResolver) ID(ctx context.Context, obj *model.IncentiveAssetWhere, data *primitive.ObjectID) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *incentiveAssetWhereResolver) Height(ctx context.Context, obj *model.IncentiveAssetWhere, data *int) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *incentiveAssetWhereResolver) Receiver(ctx context.Context, obj *model.IncentiveAssetWhere, data *string) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *incentiveAssetWhereResolver) Assets(ctx context.Context, obj *model.IncentiveAssetWhere, data []*model.IncentiveAssetWhere) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *incentiveAssetWhereResolver) Timestamp(ctx context.Context, obj *model.IncentiveAssetWhere, data *time.Time) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *incentiveWhereResolver) Amount(ctx context.Context, obj *model.IncentiveWhere, data *int) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *incentiveWhereResolver) Denom(ctx context.Context, obj *model.IncentiveWhere, data *string) error {
	panic(fmt.Errorf("not implemented"))
}

// IncentiveAssetWhere returns generated.IncentiveAssetWhereResolver implementation.
func (r *Resolver) IncentiveAssetWhere() generated.IncentiveAssetWhereResolver {
	return &incentiveAssetWhereResolver{r}
}

// IncentiveWhere returns generated.IncentiveWhereResolver implementation.
func (r *Resolver) IncentiveWhere() generated.IncentiveWhereResolver {
	return &incentiveWhereResolver{r}
}

type incentiveAssetWhereResolver struct{ *Resolver }
type incentiveWhereResolver struct{ *Resolver }
