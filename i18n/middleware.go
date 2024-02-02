package i18n

import (
	"context"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type SettingStorage interface {
	GetLanguage() string
}

type ctxKey struct{}

func Middleware(bundle *i18n.Bundle, settingStorage SettingStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept-Language")
			localizer := i18n.NewLocalizer(bundle, settingStorage.GetLanguage(), accept)

			ctx := context.WithValue(r.Context(), ctxKey{}, localizer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getLocalizerFromContext(ctx context.Context) *i18n.Localizer {
	return ctx.Value(ctxKey{}).(*i18n.Localizer)
}
