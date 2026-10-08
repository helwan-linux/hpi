package i18n

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
)

//go:embed en/strings.json
var englishData []byte

//go:embed ar/strings.json
var arabicData []byte

//go:embed es/strings.json
var spanishData []byte

//go:embed zh/strings.json
var chineseData []byte

type Translator struct {
	language string
	strings  map[string]string
}

func New() *Translator {
	language := detectLanguage()

	data := englishData

	switch language {
	case "ar":
		data = arabicData
	case "es":
		data = spanishData
	case "zh":
		data = chineseData
	}

	translations := make(map[string]string)

	if err := json.Unmarshal(data, &translations); err != nil {
		// English is the built-in fallback.
		_ = json.Unmarshal(englishData, &translations)
		language = "en"
	}

	return &Translator{
		language: language,
		strings:  translations,
	}
}

func (t *Translator) Language() string {
	if t == nil || t.language == "" {
		return "en"
	}

	return t.language
}

func (t *Translator) Get(key string) string {
	if t == nil {
		return key
	}

	value, ok := t.strings[key]
	if ok {
		return value
	}

	// Missing translations fall back to the key itself.
	return key
}

func detectLanguage() string {
	locale := os.Getenv("LC_ALL")

	if locale == "" {
		locale = os.Getenv("LC_MESSAGES")
	}

	if locale == "" {
		locale = os.Getenv("LANG")
	}

	locale = strings.ToLower(locale)

	switch {
	case strings.HasPrefix(locale, "ar"):
		return "ar"

	case strings.HasPrefix(locale, "es"):
		return "es"

	case strings.HasPrefix(locale, "zh"):
		return "zh"

	default:
		return "en"
	}
}
