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
			c.create,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/{id}",
			c.getByID,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts",
			c.list,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/my",
			c.getMyAdverts,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/count",
			c.count,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodPatch,
			"/adverts/{id}",
			c.patch,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/{id}/approve",
			c.approve,
			jwt,
			corehttp.Role(roles.Admin),
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/{id}/reject",
			c.reject,
			jwt,
			corehttp.Role(roles.Admin),
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/{id}/archive",
			c.archive,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodDelete,
			"/adverts/{id}",
			c.delete,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/images",
			c.createImages,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/images/{id}",
			c.getImageByID,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodDelete,
			"/adverts/images/{id}",
			c.deleteImage,
			jwt,
		),

		corehttp.NewRoute(
			http.MethodPost,
			"/adverts/favourites/{id}",
			c.addToFavourites,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/favourites",
			c.listFavourites,
			jwt,
		),
		corehttp.NewRoute(
			http.MethodGet,
			"/adverts/favourites/count",
			c.countFavourites,
			jwt,
		),
	}
}
