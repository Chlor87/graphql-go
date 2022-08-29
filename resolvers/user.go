package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/Chlor87/graphql/domain"
	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/util"
)

type userResolver struct {
	model.User
	*domain.Domain
}

func (r *userResolver) WithData(d *model.User) *userResolver {
	return &userResolver{*d, r.Domain}
}

func (r *userResolver) GetUser(
	ctx context.Context,
	args *struct{ ID model.ID },
) (res *userResolver, err error) {
	user, err := r.Domain.User.Get(args.ID)
	if err != nil {
		return
	}
	return r.WithData(user), nil
}

func (r *userResolver) GetUsers() (res []*userResolver, err error) {
	users, err := r.Domain.User.List()
	if err != nil {
		return
	}
	return util.Map(r.WithData, users), nil
}

func (r *userResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.User.CreatedAt}
}

func (r *userResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.User.UpdatedAt}
}

func (r *userResolver) DeletedAt() *graphql.Time {
	if !r.User.DeletedAt.Valid {
		return nil
	}
	return &graphql.Time{Time: r.User.DeletedAt.Time}
}

func (r *userResolver) CreateUser(
	ctx context.Context,
	args *struct {
		Input *model.User
	}) (*userResolver, error) {
	input := args.Input
	newUser := model.User{Name: input.Name, Email: input.Email}
	err := r.Domain.User.Create(&newUser)
	return r.WithData(&newUser), err
}
