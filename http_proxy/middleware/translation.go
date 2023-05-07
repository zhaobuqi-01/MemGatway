package middleware

import (
	"gateway/utils"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

// TranslationMiddleware sets up the translation for Gin framework.
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化
		en := en.New()
		zh := zh.New()
		uni := ut.New(zh, zh, en)
		val := validator.New()

		// 设置翻译器
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		// 注册验证方法和翻译
		switch locale {
		case "en":
			registerTranslationsForEnglish(val, trans)
		default:
			registerTranslationsForChinese(val, trans)
		}

		c.Set(utils.TranslatorKey, trans)
		c.Set(utils.ValidatorKey, val)
		c.Next()
	}
}

func registerTranslationsForEnglish(val *validator.Validate, trans ut.Translator) {
	en_translations.RegisterDefaultTranslations(val, trans)
	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("en_comment")
	})
	registerCustomValidations(val)
	registerCustomTranslations(val, trans)
}

func registerTranslationsForChinese(val *validator.Validate, trans ut.Translator) {
	zh_translations.RegisterDefaultTranslations(val, trans)
	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})
	registerCustomValidations(val)
	registerCustomTranslations(val, trans)
}

func registerCustomValidations(val *validator.Validate) {
	val.RegisterValidation("valid_username", validUsername)
	val.RegisterValidation("valid_service_name", validServiceName)
	val.RegisterValidation("valid_rule", validRule)
	val.RegisterValidation("valid_url_rewrite", validURLRewrite)
	val.RegisterValidation("valid_header_transfor", validHeaderTransfor)
	val.RegisterValidation("valid_ipportlist", validIPPortList)
	val.RegisterValidation("valid_iplist", validIPList)
	val.RegisterValidation("valid_weightlist", validWeightList)
}

func registerCustomTranslations(val *validator.Validate, trans ut.Translator) {
	translations := []struct {
		tag           string
		registerFn    validator.RegisterTranslationsFunc
		translationFn func(ut.Translator, validator.FieldError) string
	}{
		{"valid_username", registerUsernameTranslation, translateUsername},
		{"valid_service_name", registerServiceNameTranslation, translateServiceName},
		{"valid_rule", registerRuleTranslation, translateRule},
		{"valid_url_rewrite", registerURLRewriteTranslation, translateURLRewrite},
		{"valid_header_transfor", registerHeaderTransforTranslation, translateHeaderTransfor},
		{"valid_ipportlist", registerIPPortListTranslation, translateIPPortList},
		{"valid_iplist", registerIPListTranslation, translateIPList},
		{"valid_weightlist", registerWeightListTranslation, translateWeightList},
	}

	for _, t := range translations {
		val.RegisterTranslation(t.tag, trans, t.registerFn, t.translationFn)
	}
}

func validUsername(fl validator.FieldLevel) bool {
	return fl.Field().String() == "admin"
}

func validServiceName(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
	return matched
}

func validRule(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
	return matched
}

func validURLRewrite(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	for _, ms := range strings.Split(fl.Field().String(), ",") {
		if len(strings.Split(ms, " ")) != 2 {
			return false
		}
	}
	return true
}

func validHeaderTransfor(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	for _, ms := range strings.Split(fl.Field().String(), ",") {
		if len(strings.Split(ms, " ")) != 3 {
			return false
		}
	}
	return true
}

func validIPPortList(fl validator.FieldLevel) bool {
	ipPortPattern, _ := regexp.Compile(`^\S+:\d+$`)
	for _, ms := range strings.Split(fl.Field().String(), ",") {
		if !ipPortPattern.Match([]byte(ms)) {
			return false
		}
	}
	return true
}

func validIPList(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	ipPattern, _ := regexp.Compile(`\S+`)
	for _, item := range strings.Split(fl.Field().String(), ",") {
		if !ipPattern.Match([]byte(item)) {
			return false
		}
	}
	return true
}

func validWeightList(fl validator.FieldLevel) bool {
	weightPattern, _ := regexp.Compile(`^\d+$`)
	for _, ms := range strings.Split(fl.Field().String(), ",") {
		if !weightPattern.Match([]byte(ms)) {
			return false
		}
	}
	return true
}

// Register translation functions
func registerUsernameTranslation(ut ut.Translator) error {
	return ut.Add("valid_username", "{0} 填写不正确哦", true)
}

func registerServiceNameTranslation(ut ut.Translator) error {
	return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
}

func registerRuleTranslation(ut ut.Translator) error {
	return ut.Add("valid_rule", "{0} 必须是非空字符", true)
}

func registerURLRewriteTranslation(ut ut.Translator) error {
	return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
}

func registerHeaderTransforTranslation(ut ut.Translator) error {
	return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
}

func registerIPPortListTranslation(ut ut.Translator) error {
	return ut.Add("valid_ipportlist", "{0} 不符合输入格式", true)
}

func registerIPListTranslation(ut ut.Translator) error {
	return ut.Add("valid_iplist", "{0} 不符合输入格式", true)
}

func registerWeightListTranslation(ut ut.Translator) error {
	return ut.Add("valid_weightlist", "{0} 不符合输入格式", true)
}

// Translate error functions
func translateUsername(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_username", fe.Field())
	return t
}

func translateServiceName(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_service_name", fe.Field())
	return t
}

func translateRule(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_rule", fe.Field())
	return t
}

func translateURLRewrite(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_url_rewrite", fe.Field())
	return t
}

func translateHeaderTransfor(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_header_transfor", fe.Field())
	return t
}

func translateIPPortList(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_ipportlist", fe.Field())
	return t
}

func translateIPList(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_iplist", fe.Field())
	return t
}

func translateWeightList(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("valid_weightlist", fe.Field())
	return t
}
