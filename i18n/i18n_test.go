package i18n

import (
	"fmt"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestNewBundle(t *testing.T) {
	b := NewBundle()
	test, err := i18n.NewLocalizer(b, "vi").Localize(&i18n.LocalizeConfig{
		MessageID: "Test.HelloWorld",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(test)
}
