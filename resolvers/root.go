package resolvers

import (
	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
)

type Root struct {
	// DB       *gorm.DB
	// TodoRepo *repo.Repo[model.Todo]
	*todoResolver
	*userResolver
}

func NewRoot(
	todoRepo *repo.Repo[model.Todo],
	userRepo *repo.Repo[model.User],
) *Root {
	return &Root{
		&todoResolver{Repo: todoRepo, UserRepo: userRepo},
		&userResolver{Repo: userRepo},
	}
}
