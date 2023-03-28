package middleware

import (
	"gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
)

// TranslationMiddleware sets up the translation for Gin framework.
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set supported languages
		en := en.New()
		zh := zh.New()

		// Set up the internationalization translator
		uni := ut.New(zh, zh, en)
		val := validator.New()

		// Get the translator instance based on the "locale" query parameter
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		// Register the translator with the validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			// Custom validation method
			val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})

			// Custom validation translation
			val.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
				return ut.Add("valid_username", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_username", fe.Field())
				return t
			})
		}
		c.Set(utils.TranslatorKey, trans)
		c.Set(utils.ValidatorKey, val)
		c.Next()
	}
}
