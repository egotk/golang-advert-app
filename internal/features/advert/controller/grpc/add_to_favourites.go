package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) AddToFavourites(ctx context.Context, request *advertpb.AddToFavouritesRequest) (*emptypb.Empty, error) {
	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.AddToFavouritesDTO{
		AdvertID: request.AdvertId,
		UserID:   userID,
		UserRole: userRole,
	}

	if err := c.useCase.AddToFavourites(ctx, dto); err != nil {
		return nil, fmt.Errorf("failed to add advert with id = %d to favourites: %w", request.AdvertId, err)
	}

	return &emptypb.Empty{}, nil
}
