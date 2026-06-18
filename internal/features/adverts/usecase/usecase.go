package advertusecase

import (
	"context"

	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

type UseCase struct {
	repo repo
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
}

func New(repo repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
