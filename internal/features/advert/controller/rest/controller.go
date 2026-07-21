package advertrest

import (
	"context"
	"io"
	"net/http"

	corehttp "github.com/egotk/golang-advert-app/internal/core/http"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

type Controller struct {
	useCase useCase
}

//go:generate mockgen -source=controller.go -destination=mock_usecase_test.go -package=advertrest_test
type useCase interface {
	Create(ctx context.Context, dto advertusecase.CreateDTO) (advertentity.Advert, error)
	GetByID(ctx context.Context, dto advertusecase.GetByIDDTO) (advertentity.Advert, error)
	List(ctx context.Context, dto advertusecase.ListDTO) (int64, []advertentity.Advert, error)
	Patch(ctx context.Context, dto advertusecase.PatchDTO) (advertentity.Advert, error)
	Approve(ctx context.Context, id int64) (advertentity.Advert, error)
	Reject(ctx context.Context, id int64) (advertentity.Advert, error)
	Archive(ctx context.Context, dto advertusecase.ArchiveDTO) (advertentity.Advert, error)
	Delete(ctx context.Context, dto advertusecase.DeleteDTO) error
	Count(ctx context.Context, dto advertusecase.CountDTO) (int64, error)

	CreateImages(ctx context.Context, dto advertusecase.CreateImagesDTO) (_ []advertentity.AdvertImage, err error)
	GetImageByID(ctx context.Context, dto advertusecase.GetImageDTO) (io.ReadCloser, advertentity.AdvertImage, error)
	DeleteImage(ctx context.Context, dto advertusecase.DeleteImageDTO) error

	AddToFavourites(ctx context.Context, dto advertusecase.AddToFavouritesDTO) error
	ListFavourites(ctx context.Context, dto advertusecase.ListDTO) (int64, []advertentity.Advert, error)
	CountFavourites(ctx context.Context, dto advertusecase.CountDTO) (int64, error)
}

func New(
	useCase useCase,
) *Controller {
	return &Controller{
		useCase: useCase,
	}
}

func (c *Controller) Routes(jwtService corehttp.JWTService) []corehttp.Route {
	jwt := corehttp.JWToken(jwtService)

	return []corehttp.Route{
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts",
			c.Create,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/{id}",
			c.GetByID,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts",
			c.List,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/my",
			c.GetMyAdverts,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/count",
			c.Count,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodPatch,
			"/adverts/{id}",
			c.Patch,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/{id}/approve",
			c.Approve,
			jwt,
			corehttp.Role(roles.Admin),
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/{id}/reject",
			c.Reject,
			jwt,
			corehttp.Role(roles.Admin),
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/{id}/archive",
			c.Archive,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodDelete,
			"/adverts/{id}",
			c.Delete,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/images",
			c.CreateImages,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/images/{id}",
			c.GetImageByID,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodDelete,
			"/adverts/images/{id}",
			c.DeleteImage,
			jwt,
		),

		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/favourites/{id}",
			c.AddToFavourites,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/favourites",
			c.ListFavourites,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/favourites/count",
			c.CountFavourites,
			jwt,
		),
	}
}
