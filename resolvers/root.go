package resolvers

import (
	"context"

	"gorm.io/gorm"

	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
	"github.com/Chlor87/graphql/util"
)

type Root struct {
	DB       *gorm.DB
	TodoRepo *repo.Repo[model.Todo]
}

// TODO: implement todo resolver in a separate file
type todoResolver struct {
	*model.Todo
}

func newTodoResolver(u *model.Todo) *todoResolver {
	return &todoResolver{u}
}

func (t *todoResolver) ID() model.ID {
	return t.Todo.ID
}

func (r *todoResolver) Text() string {
	return r.Todo.Text
}

func (r *Root) Todo(ctx context.Context, args struct{ ID model.ID }) (
	*todoResolver, error) {
	res, err := r.TodoRepo.Get(args.ID)
	if err != nil {
		return nil, err
	}
	return newTodoResolver(res), nil
}

func (r *Root) Todos(ctx context.Context) (res []*todoResolver, err error) {
	tmp, err := r.TodoRepo.List()
	if err != nil {
		return
	}
	res = util.Map(newTodoResolver, tmp)
	return
}

func (r *Root) CreateTodo(ctx context.Context, args struct{ Input struct{ Text string } }) (*todoResolver, error) {
	t := model.Todo{Text: args.Input.Text, UpdatedByID: 4}
	err := r.TodoRepo.Create(&t)
	if err != nil {
		return nil, err
	}
	return newTodoResolver(&t), nil
}
