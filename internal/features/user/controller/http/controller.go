package userhttp

import (
	"context"
	"net/http"

	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userdto "github.com/egotk/golang-advert-app/internal/features/user/usecase/dto"
)

type Controller struct {
	useCase useCase
}

type useCase interface {
	Create(
		ctx context.Context,
		dto userdto.Create,
	) (userentity.User, error)

	List(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]userentity.User, error)

	GetByID(
		ctx context.Context,
		id int,
	) (userentity.User, error)
}

func New(useCase useCase) *Controller {
	return &Controller{
		useCase: useCase,
	}
}

func (c *Controller) Routes() []corehttp.Route {
	return []corehttp.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: c.create,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: c.list,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: c.getByID,
		},
	}
}
