package domain

import (
	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
	"gorm.io/gorm"
)

type Domain struct {
	User *repo.Repo[model.User]
	Todo *repo.Repo[model.Todo]
}

func New(db *gorm.DB) (res *Domain, err error) {
	userRepo, err := repo.New[model.User](db)
	if err != nil {
		return
	}

	todoRepo, err := repo.New[model.Todo](db)
	if err != nil {
		return
	}

	res = &Domain{userRepo, todoRepo}
	return
}
