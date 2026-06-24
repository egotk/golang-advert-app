package advertusecase

import (
	"context"
	"io"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

type UseCase struct {
	repo    repo
	storage storage
}

type repo interface {
	Create(
		ctx context.Context,
		advert *advertentity.Advert,
	) error

	GetByID(
		ctx context.Context,
		id int,
	) (advertentity.Advert, error)

	List(
		ctx context.Context,
		limit *int,
		offset *int,
		filter advertentity.Filter,
	) ([]advertentity.Advert, error)

	Patch(
		ctx context.Context,
		advert *advertentity.Advert,
	) error

	ChangeStatus(
		ctx context.Context,
		id int,
		oldStatus advertentity.Status,
		newStatus advertentity.Status,
	) (advertentity.Advert, error)

	IncrementViewsCount(
		ctx context.Context,
		id int,
	) error

	DeleteByID(
		ctx context.Context,
		id int,
	) error

	Count(
		ctx context.Context,
		filter advertentity.Filter,
	) (int, error)

	CreateImages(
		ctx context.Context,
		advertID int,
		images []advertentity.AdvertImage,
	) error

	GetImageByID(
		ctx context.Context,
		imageID int,
	) (int, advertentity.AdvertImage, error)

	ListImagesByAdvertID(
		ctx context.Context,
		advertID int,
	) ([]advertentity.AdvertImage, error)

	ListImagesByAdvertIDs(
		ctx context.Context,
		ids []int,
	) (map[int][]advertentity.AdvertImage, error)

	DeleteImageByID(
		ctx context.Context,
		advertID int,
	) error
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
	repo repo,
	storage storage,
) *UseCase {
	return &UseCase{
		repo:    repo,
		storage: storage,
	}
}
