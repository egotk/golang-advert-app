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

func RegisterValidations(validations ...Validation) error {
	for _, v := range validations {
		if err := instance.RegisterValidation(v.Tag, v.Fn); err != nil {
			return fmt.Errorf("register validation %s: %w", v.Tag, err)
		}
	}

	return nil
}

func Instance() *validator.Validate {
	return instance
}
