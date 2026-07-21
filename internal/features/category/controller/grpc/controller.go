package categorygrpc

import (
	"context"

	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
)

type Controller struct {
	categorypb.UnimplementedCategoryServer
	useCase useCase
}

//go:generate mockgen -source=controller.go -destination=mock_usecase_test.go -package=categorygrpc_test
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
