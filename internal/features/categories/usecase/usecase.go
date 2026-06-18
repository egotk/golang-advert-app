package categoryusecase

import (
	"context"

	categoryentity "github.com/egotk/golang-advert-app/internal/features/categories/entity"
)

type UseCase struct {
	repo repo
}

type repo interface {
	Create(
		ctx context.Context,
		category *categoryentity.Category,
	) error

	List(ctx context.Context) ([]categoryentity.Category, error)

	GetByID(
		ctx context.Context,
		id int,
	) (categoryentity.Category, error)

	Patch(
		ctx context.Context,
		category categoryentity.Category,
	) error

	DeleteByID(
		ctx context.Context,
		id int,
	) error
}

func New(repo repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
