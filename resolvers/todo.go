package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
	"github.com/Chlor87/graphql/util"
)

type todoResolver struct {
	model.Todo
	*repo.Repo[model.Todo]
	UserRepo *repo.Repo[model.User]
}

func (r *todoResolver) WithData(d *model.Todo) *todoResolver {
	return &todoResolver{*d, r.Repo, r.UserRepo}
}

func (u *todoResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: u.Todo.CreatedAt}
}

func (u *todoResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: u.Todo.UpdatedAt}
}

func (r *todoResolver) GetTodo(
	ctx context.Context, args struct{ ID model.ID },
) (*todoResolver, error) {
	res, err := r.Get(args.ID)
	if err != nil {
		return nil, err
	}
	return r.WithData(res), nil
}

func (r *todoResolver) GetTodos(
	ctx context.Context,
) (res []*todoResolver, err error) {
	tmp, err := r.List()
	if err != nil {
		return
	}
	res = util.Map(r.WithData, tmp)
	return
}

func (r *todoResolver) CreateTodo(
	ctx context.Context,
	args struct {
		Input struct{ Text string }
	},
) (*todoResolver, error) {
	t := model.Todo{Text: args.Input.Text, UpdatedByID: 4}
	err := r.Create(&t)
	if err != nil {
		return nil, err
	}
	return r.WithData(&t), nil
}

func (r *todoResolver) UpdatedBy(ctx context.Context) (*userResolver, error) {
	u, err := r.UserRepo.Get(model.ID(r.Todo.UpdatedByID))
	if err != nil {
		return nil, err
	}
	return &userResolver{*u, r.UserRepo}, nil
}
