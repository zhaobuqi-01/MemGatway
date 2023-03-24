package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

// DefaultGetValidParams 根据请求参数对结构体进行绑定并验证，返回错误信息
func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	// 使用 ShouldBind 绑定请求参数到结构体
	if err := c.ShouldBind(params); err != nil {
		return err
	}

	// 获取验证器
	valid, err := GetValidator(c)
	if err != nil {
		return err
	}

	// 获取翻译器
	trans, err := GetTranslation(c)
	if err != nil {
		return err
	}

	// 使用验证器进行结构体校验
	err = valid.Struct(params)
	if err != nil {
		// 将错误信息翻译为英文并返回
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			// 调用翻译器进行错误信息翻译
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}

	return nil
}

// GetValidator 从上下文中获取 validator 实例并返回，若未设置则返回错误
func GetValidator(c *gin.Context) (*validator.Validate, error) {
	val, ok := c.Get(ValidatorKey)
	if !ok {
		return nil, errors.New("Validator instance not found in context")
	}

	validator, ok := val.(*validator.Validate)
	if !ok {
		return nil, errors.New("Failed to get validator instance from context")
	}

	return validator, nil
}

// GetTranslation 从上下文中获取翻译器实例并返回，若未设置则返回错误
func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(TranslatorKey)
	if !ok {
		return nil, errors.New("Translator instance not found in context")
	}

	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, errors.New("Failed to get translator instance from context")
	}

	return translator, nil
}
