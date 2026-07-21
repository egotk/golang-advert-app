package advertusecase

import (
	"context"
	"io"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

type UseCase struct {
	favRepo favRepo
	repo    repo
	storage storage
}

//go:generate mockgen -source=usecase.go -destination=mock_usecase_test.go -package=advertusecase_test
type repo interface {
	Create(ctx context.Context, advert *advertentity.Advert) error
	GetByID(ctx context.Context, id int64) (advertentity.Advert, error)
	List(ctx context.Context, limit *int64, offset *int64, filter advertentity.Filter) ([]advertentity.Advert, error)
	Patch(ctx context.Context, advert *advertentity.Advert) error
	ChangeStatus(ctx context.Context, id int64, oldStatus advertentity.Status, newStatus advertentity.Status) (advertentity.Advert, error)
	IncrementViewsCount(ctx context.Context, id int64) error
	DeleteByID(ctx context.Context, id int64) error
	Count(ctx context.Context, filter advertentity.Filter) (int64, error)

	CreateImages(ctx context.Context, advertID int64, images []advertentity.AdvertImage) error
	GetImageByID(ctx context.Context, imageID int64) (int64, advertentity.AdvertImage, error)
	ListImagesByAdvertID(ctx context.Context, advertID int64) ([]advertentity.AdvertImage, error)
	ListImagesByAdvertIDs(ctx context.Context, ids []int64) (map[int64][]advertentity.AdvertImage, error)
	DeleteImageByID(ctx context.Context, advertID int64) error

	ListFavourites(ctx context.Context, userID int64, limit *int64, offset *int64, filter advertentity.Filter) ([]advertentity.Advert, error)
	CountFavourites(ctx context.Context, userID int64, filter advertentity.Filter) (int64, error)
}

type favRepo interface {
	Add(ctx context.Context, advertID int64, userID int64) error
}

type storage interface {
	Save(
		extension string,
		reader io.Reader,
	) (string, error)

	GetByPath(path string) (io.ReadCloser, error)

	DeleteByPath(path string) error
}

func New(
	favRepo favRepo,
	repo repo,
	storage storage,
) *UseCase {
	return &UseCase{
		favRepo: favRepo,
		repo:    repo,
		storage: storage,
	}
}
