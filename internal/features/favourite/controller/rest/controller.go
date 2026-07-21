package favrest

import (
	"context"
	"net/http"

	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
)

type Controller struct {
	useCase useCase
}

//go:generate mockgen -source=controller.go -destination=mock_usecase_test.go -package=favrest_test
type useCase interface {
	Remove(ctx context.Context, dto favusecase.RemoveDTO) error
	ListIDs(ctx context.Context, userID int64) ([]int64, error)
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
			http.MethodGet,
			"/adverts/favourites/ids",
			c.ListIDs,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodDelete,
			"/adverts/favourites/{id}",
			c.Remove,
			jwt,
		),
	}
}
