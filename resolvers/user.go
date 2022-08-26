package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
)

type userResolver struct {
	model.User
	*repo.Repo[model.User]
}

func (r *userResolver) WithData(d *model.User) *userResolver {
	return &userResolver{*d, r.Repo}
}

func (u *userResolver) GetUser(
	ctx context.Context,
	args struct{ ID model.ID },
) (*userResolver, error) {
	user, err := u.Get(args.ID)
	if err != nil {
		return nil, err
	}
	return u.WithData(user), nil
}

func (u *userResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: u.User.CreatedAt}
}

func (u *userResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: u.User.UpdatedAt}
}
