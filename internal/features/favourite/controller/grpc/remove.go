package favgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
	favpb "github.com/egotk/golang-advert-app/internal/gen/favourite"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Remove(ctx context.Context, request *favpb.RemoveRequest) (*emptypb.Empty, error) {
	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get 'UserInfo' from JWT: %w", err)
	}

	dto := favusecase.RemoveDTO{
		AdvertID: request.AdvertId,
		UserID:   userID,
	}

	if err := c.useCase.Remove(ctx, dto); err != nil {
		return nil, fmt.Errorf("failed to remove advert from favourites: %w", err)
	}

	return &emptypb.Empty{}, nil
}
