package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Chlor87/graphql/graph/generated"
	"github.com/Chlor87/graphql/middleware"
	"github.com/Chlor87/graphql/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	curr, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	t := model.Todo{
		Text:        input.Text,
		UpdatedByID: curr.ID,
	}
	err = r.TodoRepo.Create(&t)
	return &t, err
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id int, input model.NewTodo) (*model.Todo, error) {
	curr, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	t := model.Todo{
		Text:        input.Text,
		UpdatedByID: curr.ID,
	}
	return r.TodoRepo.Update(id, &t)
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	u := model.User{
		Name:  input.Name,
		Email: input.Email,
	}
	err := r.UserRepo.Create(&u)
	return &u, err
}

// Todo is the resolver for the todo field.
func (r *queryResolver) Todo(ctx context.Context, id int) (*model.Todo, error) {
	return r.TodoRepo.Get(id)
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.TodoRepo.List()
}

// MyTodos is the resolver for the myTodos field.
func (r *queryResolver) MyTodos(ctx context.Context) ([]*model.Todo, error) {
	curr, err := middleware.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	var res []*model.Todo
	err = r.TodoRepo.DB.
		Order("id asc").
		Where(&model.Todo{UpdatedByID: curr.ID}).
		Find(&res).Error
	return res, err
}

// SearchTodos is the resolver for the searchTodos field.
func (r *queryResolver) SearchTodos(ctx context.Context, input model.TodoSearch) ([]*model.Todo, error) {
	tx := r.TodoRepo.DB.Model(&model.Todo{})
	if input.CreatedAt != nil {
		if input.CreatedAt.From != nil {
			tx.Where("created_at > ?", input.CreatedAt.From)
		}
		if input.CreatedAt.To != nil {
			tx.Where("created_at < ?", input.CreatedAt.To)
		}
	}

	var res []*model.Todo
	err := tx.Find(&res).Error
	return res, err
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	return r.UserRepo.Get(id)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.UserRepo.List()
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
