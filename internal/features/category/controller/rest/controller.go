package categoryhttp

import (
	"context"
	"net/http"

	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
)

type Controller struct {
	useCase useCase
}

//go:generate mockgen -source=controller.go -destination=mock_usecase_test.go -package=categoryhttp_test
type useCase interface {
	Create(ctx context.Context, dto categoryusecase.CreateDTO) (categoryentity.Category, error)
	List(ctx context.Context) ([]categoryentity.Category, error)
	Patch(ctx context.Context, dto categoryusecase.PatchDTO) (categoryentity.Category, error)
	Delete(ctx context.Context, id int64) error
}

func New(useCase useCase) *Controller {
	return &Controller{
		useCase: useCase,
	}
}

func (c *Controller) Routes(jwtService corehttp.JWTService) []corehttp.Route {
	jwt := corehttp.JWToken(jwtService)

	return []corehttp.Route{
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/categories",
			c.Create,
			jwt,
			corehttp.Role(roles.Admin),
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/categories",
			c.List,
		),
		corehttp.NewRoute(
			http.MethodPatch,
			"/adverts/categories/{id}",
			c.Patch,
			jwt,
			corehttp.Role(roles.Admin),
		),
		corehttp.NewRoute(
			http.MethodDelete,
			"/adverts/categories/{id}",
			c.Delete,
			jwt,
			corehttp.Role(roles.Admin),
		),
	}
}
