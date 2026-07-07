package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) CountFavourites(ctx context.Context, request *advertpb.CountRequest) (*advertpb.CountResponse, error) {
	filter := filterFromRequest(request.Filter)

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.CountDTO{
		UserID:   userID,
		UserRole: userRole,
		Filter:   filter,
	}

	favCount, err := c.useCase.CountFavourites(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to count favourites: %w", err)
	}

	response := &advertpb.CountResponse{Count: favCount}

	return response, nil
}
