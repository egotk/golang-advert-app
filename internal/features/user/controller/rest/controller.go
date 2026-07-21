package userhttp

import (
	"context"
	"net/http"

	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
)

type Controller struct {
	useCase useCase
}

//go:generate mockgen -source=controller.go -destination=mock_usecase_test.go -package=userhttp_test
type useCase interface {
	CreateUser(ctx context.Context, dto userusecase.CreateDTO) (userentity.User, error)
	ListUsers(ctx context.Context, limit *int64, offset *int64) ([]userentity.User, error)
	GetUserByID(ctx context.Context, id int64) (userentity.User, error)

	Login(ctx context.Context, dto userusecase.LoginDTO) (userusecase.LoginResultDTO, error)
	Logout(ctx context.Context, dto userusecase.LogoutDTO) error
	RefreshTokens(ctx context.Context, dto userusecase.RefreshTokensDTO) (userusecase.TokensDTO, error)
}

func New(useCase useCase) *Controller {
	return &Controller{
		useCase: useCase,
	}
}

func (c *Controller) Routes(jwtService corehttp.JWTService) []corehttp.Route {
	jwt := corehttp.JWToken(jwtService)

	return []corehttp.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: c.CreateUser,
		},
		{
			Method:     http.MethodGet,
			Path:       "/users",
			Handler:    c.ListUsers,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:     http.MethodGet,
			Path:       "/users/{id}",
			Handler:    c.GetUserByID,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: c.Login,
		},
		{
			Method:     http.MethodPost,
			Path:       "/auth/logout",
			Handler:    c.Logout,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: c.RefreshTokens,
		},
	}
}
