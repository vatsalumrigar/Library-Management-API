package localization

import (
	"encoding/json"
	"fmt"
	// "net/http"

	// "github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	logs "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var Bundel *i18n.Bundle

func LoadBundel() *i18n.Bundle {

	Languages := [2]string{"en","hn"}

	Bundel = i18n.NewBundle(language.English)
	Bundel.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range Languages {
		logs.Debugln("loading file", fmt.Sprintf("localise/%v.json", lang))
		Bundel.LoadMessageFile(fmt.Sprintf("localise/%v.json", lang))
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

// func GetLanguage(c *gin.Context) bool {

// 	languageToken := c.Request.Header.Get("lan")
		
// 	if languageToken == "" {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "No language header provided"})
// 		c.Abort()
// 		logs.Error("No langgauge header provided")
// 		return false
// 	}

// 	return true
// }