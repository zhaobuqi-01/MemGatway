package utils

import (
	"fmt"
	"gateway/globals"
	"gateway/pkg/log"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	// Bind request parameters to struct
	if err := c.ShouldBind(params); err != nil {
		log.Error("failed to bind request parameters", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	// Get validator
	valid, err := GetValidator(c)
	if err != nil {
		log.Error("failed to get validator", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	// Get translator
	trans, err := GetTranslation(c)
	if err != nil {
		log.Error("failed to get translator", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
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
		log.Error("invalid parameters", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("invalid parameters: %s", strings.Join(sliceErrs, ";"))
	}

	return nil
}

func GetValidator(c *gin.Context) (*validator.Validate, error) {
	val, ok := c.Get(globals.ValidatorKey)
	if !ok {
		return nil, fmt.Errorf("validator instance not found in context")
	}

	validator, ok := val.(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to get validator instance from context")
	}

	return validator, nil
}

func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(globals.TranslatorKey)
	if !ok {
		return nil, fmt.Errorf("translator instance not found in context")
	}

	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, fmt.Errorf("failed to get translator instance from context")
	}

	return translator, nil
}
