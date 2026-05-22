package corevalidator

import "github.com/go-playground/validator/v10"

var instance = validator.New()

func Instance() *validator.Validate {
	return instance
}
