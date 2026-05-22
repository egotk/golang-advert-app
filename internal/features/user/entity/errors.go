package userentity

import "errors"

var (
	ErrEmptyPatch = errors.New("empty patch request")

	ErrEmailCantBeNull    = errors.New("'Email' cant be null")
	ErrInvalidEmailLen    = errors.New("invalid 'Email' len")
	ErrInvalidEmailFormat = errors.New("invalid 'Email' format")

	ErrFullNameCantBeNull = errors.New("'FullName' cant be null")
	ErrInvalidFullNameLen = errors.New("invalid 'FullName' len")

	ErrPhoneNumberCantBeNull    = errors.New("'PhoneNumber' cant be null")
	ErrInvalidPhoneNumberLen    = errors.New("invalid 'PhoneNumber' len")
	ErrInvalidPhoneNumberFormat = errors.New("invalid 'PhoneNumber' format")

	ErrPasswordCantBeNull = errors.New("'Password' cant be null")
	ErrInvalidPasswordLen = errors.New("invalid 'Password' len")

	ErrRoleCantBeNull = errors.New("'Role' cant be null")
	ErrInvalidRole    = errors.New("'Role' is invalid")
)
