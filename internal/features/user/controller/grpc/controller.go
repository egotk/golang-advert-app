package usergrpc

import (
	"context"

	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

type Controller struct {
	userpb.UnimplementedUserServer
	useCase useCase

	log *corezaplogger.Logger
}

type useCase interface {
	CreateUser(ctx context.Context, dto userusecase.CreateDTO) (userentity.User, error)
	ListUsers(ctx context.Context, limit *int64, offset *int64) ([]userentity.User, error)
	GetUserByID(ctx context.Context, id int64) (userentity.User, error)

	Login(ctx context.Context, dto userusecase.LoginDTO) (userusecase.LoginResultDTO, error)
	Logout(ctx context.Context, dto userusecase.LogoutDTO) error
	RefreshTokens(ctx context.Context, dto userusecase.RefreshTokensDTO) (userusecase.TokensDTO, error)
}

func New(
	useCase useCase,
	log *corezaplogger.Logger,
) *Controller {
	return &Controller{
		useCase: useCase,
		log:     log,
	}
}
