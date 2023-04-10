package pkg

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

const (
	ValidatorKey  = "ValidatorKey"
	TranslatorKey = "TranslatorKey"
)

func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	// Bind request parameters to struct
	if err := c.ShouldBind(params); err != nil {
		return err
	}

	// Get validator
	valid, err := GetValidator(c)
	if err != nil {
		return err
	}

	// Get translator
	trans, err := GetTranslation(c)
	if err != nil {
		return err
	}

	// Validate struct
	err = valid.Struct(params)
	if err != nil {
		// Translate error messages and return
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}

	return nil
}

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
