// TODO: nest deeper and rename to slice
package util

func Map[T any, U any](fn func(T) U, xs []T) []U {
	res := make([]U, len(xs))
	for i, x := range xs {
		res[i] = fn(x)
	}
	return res
}
