package userentity

import (
	"regexp"

	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
)

const (
	minEmailLen = 3
	maxEmailLen = 255

	minFullNameLen = 3
	maxFullNameLen = 100

	minPhoneNumberLen = 4
	maxPhoneNumberLen = 20

	minPasswordLen     = 8
	maxPasswordLen     = 64
	maxPasswordByteLen = 72 // ограничение bcrypt
)

var (
	userRoles = map[string]struct{}{
		"user":  {},
		"admin": {},
	}

	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
)

func ValidateEmail(email string) error {
	emailLen := len([]rune(email))
	if emailLen < minEmailLen || emailLen > maxEmailLen {
		return ErrInvalidEmailLen
	}

	validator := corevalidator.Instance()
	if err := validator.Var(email, "email"); err != nil {
		return ErrInvalidEmailFormat
	}

	return nil
}

func ValidateFullName(fullName string) error {
	fullNameLen := len([]rune(fullName))
	if fullNameLen < minFullNameLen || fullNameLen > maxFullNameLen {
		return ErrInvalidFullNameLen
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	phoneNumberLen := len([]rune(phoneNumber))
	if phoneNumberLen < minPhoneNumberLen || phoneNumberLen > maxPhoneNumberLen {
		return ErrInvalidPhoneNumberLen
	}

	if !phoneRegex.MatchString(phoneNumber) {
		return ErrInvalidPhoneNumberFormat
	}

	return nil
}

func ValidatePassword(password string) error {
	passwordLen := len([]rune(password))
	passwordByteLen := len(password)

	if passwordLen < minPasswordLen || passwordLen > maxPasswordLen || passwordByteLen > maxPasswordByteLen {
		return ErrInvalidPasswordLen
	}

	return nil
}

func ValidateRole(role string) error {
	_, valid := userRoles[role]
	if !valid {
		return ErrInvalidRole
	}

	return nil
}

func ValidateImagePath(imagePath string) error {
	// TODO:
	return nil
}
