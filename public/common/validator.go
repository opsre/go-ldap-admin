package common

import (
	"regexp"
	"strings"

	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	ch_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局Validate数据校验实列
var Validate *validator.Validate

// 全局翻译器
var Trans ut.Translator

var TransByLocale map[string]ut.Translator

// 初始化Validator数据校验
func InitValidate() {
	Validate = validator.New()
	TransByLocale = map[string]ut.Translator{}

	chinese := zh.New()
	english := en.New()
	japanese := ja.New()
	spanish := es.New()
	uni := ut.New(chinese, chinese, english, japanese, spanish)

	register := func(locale string, registerFn func(*validator.Validate, ut.Translator) error) {
		trans, _ := uni.GetTranslator(locale)
		_ = registerFn(Validate, trans)
		TransByLocale[locale] = trans
	}

	register("zh", ch_translations.RegisterDefaultTranslations)
	register("en", en_translations.RegisterDefaultTranslations)
	register("ja", ja_translations.RegisterDefaultTranslations)
	register("es", es_translations.RegisterDefaultTranslations)
	TransByLocale["ko"] = TransByLocale["en"]

	Trans = TransByLocale["zh"]
	_ = Validate.RegisterValidation("checkMobile", checkMobile)
	Log.Infof("初始化validator.v10数据校验器完成")
}

func TranslatorForLocale(locale string) ut.Translator {
	lang := strings.Split(i18n.NormalizeLocale(locale), "-")[0]
	if trans, ok := TransByLocale[lang]; ok {
		return trans
	}
	return Trans
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `1\d{10}`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}
