package echox

import (
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enLang "github.com/go-playground/validator/v10/translations/en"
	zhLang "github.com/go-playground/validator/v10/translations/zh"
)

var (
	v          *validator.Validate
	translator *ut.UniversalTranslator
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func initValidate() {
	v = validator.New()

	translator = ut.New(en.New(), en.New(), zh.New())
	if en, success := translator.GetTranslator("en"); success {
		enLang.RegisterDefaultTranslations(v, en)
	}
	if zh, success := translator.GetTranslator("zh"); success {
		zhLang.RegisterDefaultTranslations(v, zh)
	}
}

func i18n(lang string, errs validator.ValidationErrors) (i18n validator.ValidationErrorsTranslations) {
	sep := "_"
	if strings.Contains(lang, "-") {
		sep = "-"
	}

	splits := strings.Split(lang, sep)
	for i := 0; i < len(splits); i++ {
		if t, s := translator.FindTranslator(lang); s {
			i18n = errs.Translate(t)
			break
		} else {
			if index := strings.LastIndex(lang, sep); -1 == index {
				break
			} else {
				lang = lang[0:index]
			}
		}
	}
	if nil == i18n {
		if t, s := translator.GetTranslator("zh"); s {
			i18n = errs.Translate(t)
		}
	}

	return
}
