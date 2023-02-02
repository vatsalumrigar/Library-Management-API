package localization

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundel1 *i18n.Bundle

func LoadBundel1() *i18n.Bundle {

	Languages := [2]string{"en","hi"}


	Bundel1 = i18n.NewBundle(language.English)
	Bundel1.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range Languages {
		Bundel1.MustLoadMessageFile(fmt.Sprintf("localise/%v.json", lang))
	}

	return Bundel1
}

func GetMessage1(lang string, id string) string {

	localizer := i18n.NewLocalizer(Bundel1, lang)

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
