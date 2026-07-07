package favgrpc

import (
	"context"

	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
	favpb "github.com/egotk/golang-advert-app/internal/gen/favourite"
)

type Controller struct {
	favpb.UnimplementedFavouriteServer
	useCase useCase
}

type useCase interface {
	Remove(ctx context.Context, dto favusecase.RemoveDTO) error
	ListIDs(ctx context.Context, userID int64) ([]int64, error)
}

func New(
	useCase useCase,
) *Controller {
	return &Controller{
		useCase: useCase,
	}
}
