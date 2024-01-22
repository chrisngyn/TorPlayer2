package request

import (
	"context"
	"net/http"
)

type Context struct {
	URL string
}

type ctxKey struct{}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKey{}, Context{
			URL: r.URL.String(),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetContext(ctx context.Context) Context {
	return ctx.Value(ctxKey{}).(Context)
}
