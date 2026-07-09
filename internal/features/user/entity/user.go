package userentity

import "time"

const (
	initialVersion = 1
	initialRole    = "user"
)

type User struct {
	ID               int64
	Version          int64
	Email            string
	FullName         string
	PhoneNumber      string
	PasswordHash     string
	Role             string
	FailedLoginCount int64
	LockedUntil      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func New(
	id int64,
	version int64,
	email string,
	fullName string,
	phoneNumber string,
	passwordHash string,
	role string,
	failedLoginCount int64,
	lockedUntil *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	imagePath *string,
) User {
	return User{
		ID:               id,
		Version:          version,
		Email:            email,
		FullName:         fullName,
		PhoneNumber:      phoneNumber,
		PasswordHash:     passwordHash,
		Role:             role,
		FailedLoginCount: failedLoginCount,
		LockedUntil:      lockedUntil,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func NewInitial(
	email string,
	fullName string,
	phoneNumber string,
	passwordHash string,
) User {
	now := time.Now()

	return New(
		0,
		initialVersion,
		email,
		fullName,
		phoneNumber,
		passwordHash,
		initialRole,
		0,
		nil,
		now,
		now,
		nil,
	)
}
