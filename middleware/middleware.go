package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/Chlor87/graphql/loaders"
	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
)

type ctxKey int

const (
	userKey ctxKey = iota
	userLoaderKey
)

func fail(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Build(fns ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(fns) - 1; i >= 0; i-- {
			next = fns[i](next)
		}
		return next
	}
}

func AddUser(userRepo *repo.Repo[model.User]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := userRepo.Get(4)
			if err != nil {
				fail(w, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, u)))
		})
	}
}

func AddUserLoader(userRepo *repo.Repo[model.User]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userLoader := loaders.NewUserLoader(loaders.UserLoaderConfig{
				MaxBatch: 100,
				Wait:     1 * time.Millisecond,
				Fetch: func(keys []int) ([]*model.User, []error) {
					var tmp []*model.User
					err := userRepo.DB.Where("id in (?)", keys).Find(&tmp).Error
					if err != nil {
						return nil, []error{err}
					}

					m := make(map[int]*model.User)

					for _, u := range tmp {
						m[u.ID] = u
					}

					res := make([]*model.User, len(m))

					for i, id := range keys {
						res[i] = m[id]
					}

					return res, nil
				},
			})
			ctx := context.WithValue(r.Context(), userLoaderKey, userLoader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getCtxVal[T any](ctx context.Context, key ctxKey) (T, error) {
	v, ok := ctx.Value(key).(T)
	if !ok {
		var t T
		return t, fmt.Errorf(
			"failed to read and assert value from context %s", reflect.TypeOf(t))
	}
	return v, nil
}

func GetUser(ctx context.Context) (*model.User, error) {
	return getCtxVal[*model.User](ctx, userKey)
}

func GetUserLoader(ctx context.Context) (*loaders.UserLoader, error) {
	return getCtxVal[*loaders.UserLoader](ctx, userLoaderKey)
}
