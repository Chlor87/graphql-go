package middleware

import (
	"log"
	"net/http"
)

type ctxKey int

type MiddlewareFunc func(http.Handler) http.Handler

const (
	userKey ctxKey = iota
	userLoaderKey
)

func fail(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Build performs a left-to-right function composition
func Build(fns ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		if len(fns) == 0 {
			return next
		}
		return fns[0](Build(fns[1:]...)(next))
	}
}
