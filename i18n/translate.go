package i18n

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func T(ctx context.Context, messageID string) (string, error) {
	return getLocalizerFromContext(ctx).Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

func MustT(ctx context.Context, messageID string) string {
	s, _ := T(ctx, messageID)
	return s
}

func TWithData(ctx context.Context, messageID string, templateData, PluralCount any) (string, error) {
	return getLocalizerFromContext(ctx).Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  PluralCount,
	})
}
