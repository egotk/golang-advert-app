package advertgrpc

import (
	"context"
	"io"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

type Controller struct {
	advertpb.UnimplementedAdvertServer
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
