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

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v\n",
			r.Method, r.RemoteAddr, r.URL, time.Now().Sub(start))
	})
}

func AddUser(userRepo *repo.Repo[model.User]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := userRepo.Get(1)
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
				Fetch: func(keys []model.ID) ([]*model.User, []error) {
					var tmp []*model.User
					err := userRepo.DB.Where("id in (?)", keys).Find(&tmp).Error
					if err != nil {
						return nil, []error{err}
					}

					m := make(map[int]*model.User)

					for _, u := range tmp {
						m[int(u.ID)] = u
					}

					res := make([]*model.User, len(keys))

					for i, id := range keys {
						res[i] = m[int(id)]
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
