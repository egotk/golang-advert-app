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

//go:generate mockgen -source=usecase.go -destination=mock_usecase_test.go -package=userusecase_test
type repo interface {
	CreateUser(
		ctx context.Context,
		user *userentity.User,
	) error

	ListUsers(
		ctx context.Context,
		limit *int64,
		offset *int64,
	) ([]userentity.User, error)

	GetUserByID(
		ctx context.Context,
		id int64,
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
		userId int64,
		hash string,
	) error

	ReissueRefreshToken(
		ctx context.Context,
		userID int64,
		oldHash string,
		newToken userentity.RefreshToken,
	) error

	IncrementFailedLoginCount(
		ctx context.Context,
		id int64,
		version int64,
	) (*time.Time, error)

	ResetFailedLoginCount(
		ctx context.Context,
		id int64,
		version int64,
	) error
}

type jwtService interface {
	IssuePair(
		role string,
		userId int64,
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
