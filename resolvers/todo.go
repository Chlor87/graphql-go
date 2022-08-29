package resolvers

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"

	"github.com/Chlor87/graphql/domain"
	mw "github.com/Chlor87/graphql/middleware"
	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/util"
)

type todoResolver struct {
	model.Todo
	*domain.Domain
}

func (r *todoResolver) WithData(d *model.Todo) *todoResolver {
	return &todoResolver{*d, r.Domain}
}

func (r *todoResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.Todo.CreatedAt}
}

func (r *todoResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.Todo.UpdatedAt}
}

func (r *todoResolver) DeletedAt() *graphql.Time {
	if !r.Todo.DeletedAt.Valid {
		return nil
	}
	return &graphql.Time{Time: r.Todo.DeletedAt.Time}
}

// UpdatedBy loads todo creator using the dataloader in order to prevent
// quering the same record multiple times
func (r *todoResolver) UpdatedBy(
	ctx context.Context,
) (res *userResolver, err error) {
	loader, err := mw.GetUserLoader(ctx)
	if err != nil {
		return
	}
	user, err := loader.Load(r.Todo.UpdatedByID)
	if err != nil {
		return
	}
	res = &userResolver{*user, r.Domain}
	return
}

func (r *todoResolver) GetTodo(
	ctx context.Context,
	args *struct{ ID model.ID },
) (res *todoResolver, err error) {
	fmt.Println(args.ID)
	todo, err := r.Domain.Todo.Get(args.ID)
	if err != nil {
		return
	}
	return r.WithData(todo), nil
}

func (r *todoResolver) GetTodos(
	ctx context.Context,
) (res []*todoResolver, err error) {
	tmp, err := r.Domain.Todo.List()
	if err != nil {
		return
	}
	res = util.Map(r.WithData, tmp)
	return
}

func (r *todoResolver) CreateTodo(
	ctx context.Context,
	args *struct {
		Input *model.Todo
	},
) (res *todoResolver, err error) {
	user, err := mw.GetUser(ctx)
	if err != nil {
		return
	}
	t := model.Todo{Text: args.Input.Text, UpdatedByID: user.ID}
	err = r.Domain.Todo.Create(&t)
	if err != nil {
		return
	}
	return r.WithData(&t), nil
}

func (r *todoResolver) UpdateTodo(
	ctx context.Context, args *struct {
		ID    model.ID
		Input *model.Todo
	}) (res *todoResolver, err error) {
	user, err := mw.GetUser(ctx)
	if err != nil {
		return
	}
	args.Input.UpdatedByID = user.ID
	todo, err := r.Domain.Todo.Update(args.ID, args.Input)
	if err != nil {
		return
	}
	return r.WithData(todo), nil
}

func (r *todoResolver) DeleteTodo(
	ctx context.Context,
	args *struct{ ID model.ID },
) (res *todoResolver, err error) {
	todo, err := r.Domain.Todo.Delete(args.ID)
	if err != nil {
		return
	}
	return r.WithData(todo), nil
}
