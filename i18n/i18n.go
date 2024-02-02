package i18n

import (
	"embed"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locale.*.yml
var localeFS embed.FS

func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	dirEntries, err := localeFS.ReadDir(".")
	if err != nil {
		panic(fmt.Errorf("read dir: %w", err))
	}

	for _, dirEntry := range dirEntries {
		_, err := bundle.LoadMessageFileFS(localeFS, dirEntry.Name())
		if err != nil {
			panic(fmt.Errorf("load message file %s: %w", dirEntry.Name(), err))
		}
	}

	return bundle
}
