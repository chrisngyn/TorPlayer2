package setting

import (
	"context"
	"net/http"
)

type ctxKey struct{}

func Middleware(storage *Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ctxKey{}, storage))
			next.ServeHTTP(w, r)
		})
	}
}

func GetSettingFromContext(ctx context.Context) Settings {
	return ctx.Value(ctxKey{}).(*Storage).GetSetting()
}
