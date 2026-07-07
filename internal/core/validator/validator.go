package corevalidator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var instance = validator.New()

type Validation struct {
	Tag string
	Fn  validator.Func
}

type StructValidation struct {
	Fn    validator.StructLevelFunc
	Types []any
}

func RegisterValidations(
	validations []Validation,
	structValidations []StructValidation,
) error {
	for _, v := range validations {
		if err := instance.RegisterValidation(v.Tag, v.Fn); err != nil {
			return fmt.Errorf("register validation %s: %w", v.Tag, err)
		}
	}

	for _, v := range structValidations {
		instance.RegisterStructValidation(v.Fn, v.Types...)
	}

	return nil
}

func Instance() *validator.Validate {
	return instance
}
