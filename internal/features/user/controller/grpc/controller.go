package usergrpc

import (
	"context"

	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

type Controller struct {
	userpb.UnimplementedUserServer
	useCase useCase
}

//go:generate mockgen -source=controller.go -destination=mock_usecase_test.go -package=usergrpc_test
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
) *Controller {
	return &Controller{
		useCase: useCase,
	}
}
