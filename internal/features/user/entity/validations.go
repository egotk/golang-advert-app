package userentity

import (
	"regexp"

	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	"github.com/go-playground/validator/v10"
)

const maxPasswordByteLen = 72 // ограничение bcrypt

var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func Validations() []corevalidator.Validation {
	return []corevalidator.Validation{
		{
			Tag: "bcrypt_password_byte_len",
			Fn: func(fl validator.FieldLevel) bool {
				return len(fl.Field().String()) <= maxPasswordByteLen
			},
		},
		{
			Tag: "phone_regex",
			Fn: func(fl validator.FieldLevel) bool {
				return phoneRegex.MatchString(fl.Field().String())
			},
		},
	}
}
