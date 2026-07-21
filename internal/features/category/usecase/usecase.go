package categoryusecase

import (
	"context"

	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

type UseCase struct {
	repo repo
}

//go:generate mockgen -source=usecase.go -destination=mock_usecase_test.go -package=categoryusecase_test
type repo interface {
	Create(
		ctx context.Context,
		category *categoryentity.Category,
	) error

	List(ctx context.Context) ([]categoryentity.Category, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (categoryentity.Category, error)

	Patch(
		ctx context.Context,
		category categoryentity.Category,
	) error

	DeleteByID(
		ctx context.Context,
		id int64,
	) error
}

func New(repo repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
