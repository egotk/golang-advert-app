package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) ListFavourites(ctx context.Context, request *advertpb.ListRequest) (*advertpb.AdvertsResponse, error) {
	filter := filterFromRequest(request.Filter)

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.ListDTO{
		UserID:   userID,
		UserRole: userRole,
		Limit:    request.Limit,
		Offset:   request.Offset,
		Filter:   filter,
	}

	count, favs, err := c.useCase.ListFavourites(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to list favourite adverts: %w", err)
	}

	response := advertsToResponse(favs, count)

	return response, nil
}
