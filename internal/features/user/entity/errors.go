package userentity

import (
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
)

var (
	ErrEmailCantBeNull    = fmt.Errorf("'Email' cant be null: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidEmailLen    = fmt.Errorf("invalid 'Email' len: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidEmailFormat = fmt.Errorf("invalid 'Email' format: %w", coreerrors.ErrInvalidArgument)

	ErrFullNameCantBeNull = fmt.Errorf("'FullName' cant be null: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidFullNameLen = fmt.Errorf("invalid 'FullName' len: %w", coreerrors.ErrInvalidArgument)

	ErrPhoneNumberCantBeNull    = fmt.Errorf("'PhoneNumber' cant be null: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidPhoneNumberLen    = fmt.Errorf("invalid 'PhoneNumber' len: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidPhoneNumberFormat = fmt.Errorf("invalid 'PhoneNumber' format: %w", coreerrors.ErrInvalidArgument)

	ErrPasswordCantBeNull = fmt.Errorf("'Password' cant be null: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidPasswordLen = fmt.Errorf("invalid 'Password' len: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidPassword    = fmt.Errorf("invalid password: %w", coreerrors.ErrUnauthorized)

	ErrRoleCantBeNull = fmt.Errorf("'Role' cant be null: %w", coreerrors.ErrInvalidArgument)
	ErrInvalidRole    = fmt.Errorf("'Role' is invalid: %w", coreerrors.ErrInvalidArgument)

	ErrUserVersionConflict = fmt.Errorf("user version conflict: %w", coreerrors.ErrConflict)
	ErrUserIsLocked        = fmt.Errorf("user is locked: %w", coreerrors.ErrUnauthorized)

	ErrInvalidImageFmt = fmt.Errorf("invalid image format: %w", coreerrors.ErrInvalidArgument)

	ErrRefreshTokenExpired = fmt.Errorf("refresh token expired: %w", coreerrors.ErrUnauthorized)
	ErrTokenNotFound       = fmt.Errorf("token not found: %w", coreerrors.ErrUnauthorized)
)
