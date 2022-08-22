package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/Chlor87/graphql/graph/generated"
	"github.com/Chlor87/graphql/middleware"
	"github.com/Chlor87/graphql/model"
)

// DeletedAt is the resolver for the deletedAt field.
func (r *todoResolver) DeletedAt(ctx context.Context, obj *model.Todo) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: DeletedAt - deletedAt"))
}

// UpdatedBy is the resolver for the updatedBy field.
func (r *todoResolver) UpdatedBy(ctx context.Context, obj *model.Todo) (*model.User, error) {
	loader, err := middleware.GetUserLoader(ctx)
	if err != nil {
		return nil, err
	}
	return loader.Load(obj.UpdatedByID)
}

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type todoResolver struct{ *Resolver }
