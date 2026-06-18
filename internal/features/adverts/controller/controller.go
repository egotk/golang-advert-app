package adverthttp

import (
	"context"
	"net/http"

	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/adverts/usecase"
)

type Controller struct {
	useCase useCase
}

type useCase interface {
	Create(
		ctx context.Context,
		dto advertusecase.CreateDTO,
	) (advertentity.Advert, error)

	GetByID(
		ctx context.Context,
		dto advertusecase.GetByIDDTO,
	) (advertentity.Advert, error)

	List(
		ctx context.Context,
		dto advertusecase.ListDTO,
	) (int, []advertentity.Advert, error)

	Patch(
		ctx context.Context,
		dto advertusecase.PatchDTO,
	) (advertentity.Advert, error)

	Approve(
		ctx context.Context,
		id int,
	) (advertentity.Advert, error)

	Reject(
		ctx context.Context,
		id int,
	) (advertentity.Advert, error)

	Archive(
		ctx context.Context,
		dto advertusecase.ArchiveDTO,
	) (advertentity.Advert, error)

	Delete(
		ctx context.Context,
		dto advertusecase.DeleteDTO,
	) error

	Count(
		ctx context.Context,
		dto advertusecase.CountDTO,
	) (int, error)
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
			Method:     http.MethodPost,
			Path:       "/adverts",
			Handler:    c.create,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:     http.MethodGet,
			Path:       "/adverts/{id}",
			Handler:    c.getByID,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:     http.MethodGet,
			Path:       "/adverts",
			Handler:    c.list,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:     http.MethodGet,
			Path:       "/adverts/count",
			Handler:    c.count,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/adverts/{id}",
			Handler:    c.patch,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:  http.MethodPost,
			Path:    "/adverts/{id}/approve",
			Handler: c.approve,
			Middleware: []corehttp.Middleware{
				jwt,
				corehttp.Role(roles.Admin),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/adverts/{id}/reject",
			Handler: c.reject,
			Middleware: []corehttp.Middleware{
				jwt,
				corehttp.Role(roles.Admin),
			},
		},
		{
			Method:     http.MethodPost,
			Path:       "/adverts/{id}/archive",
			Handler:    c.archive,
			Middleware: []corehttp.Middleware{jwt},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/adverts/{id}",
			Handler:    c.delete,
			Middleware: []corehttp.Middleware{jwt},
		},
	}
}
