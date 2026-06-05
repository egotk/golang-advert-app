package userusecase

import (
	"context"
	"time"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

type UseCase struct {
	repo       repo
	jwtService jwtService
}

type repo interface {
	CreateUser(
		ctx context.Context,
		user *userentity.User,
	) error

	ListUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]userentity.User, error)

	GetUserByID(
		ctx context.Context,
		id int,
	) (userentity.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (userentity.User, error)

	CreateRefreshToken(
		ctx context.Context,
		token userentity.RefreshToken,
	) error

	GetRefreshTokenByHash(
		ctx context.Context,
		hash string,
	) (userentity.RefreshToken, error)

	DeleteRefreshTokenByHash(
		ctx context.Context,
		userId int,
		hash string,
	) error

	ReissueRefreshToken(
		ctx context.Context,
		userID int,
		oldHash string,
		newToken userentity.RefreshToken,
	) error

	IncrementFailedLoginCount(
		ctx context.Context,
		id int,
		version int,
	) (*time.Time, error)

	ResetFailedLoginCount(
		ctx context.Context,
		id int,
		version int,
	) error
}

type jwtService interface {
	IssuePair(
		role string,
		userId int,
	) (corejwt.Pair, error)
}

func New(
	repo repo,
	jwtService jwtService,
) *UseCase {
	return &UseCase{
		repo:       repo,
		jwtService: jwtService,
	}
}
