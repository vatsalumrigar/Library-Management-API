package localization

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundel *i18n.Bundle

func LoadBundel() *i18n.Bundle {

	Languages := [2]string{"en","hi"}


	Bundel = i18n.NewBundle(language.English)
	Bundel.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range Languages {
		Bundel.MustLoadMessageFile(fmt.Sprintf("localise/%v.json", lang))
	}

	return Bundel
}

func GetMessage(lang string, id string) string {

	localizer := i18n.NewLocalizer(Bundel, lang)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
	
	if err != nil || message == "" {
		message = id
	}

	return message
}
