//go:generate go run github.com/99designs/gqlgen generate
package resolvers

import (
	"gorm.io/gorm"

	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB       *gorm.DB
	TodoRepo *repo.Repo[model.Todo]
	UserRepo *repo.Repo[model.User]
}
