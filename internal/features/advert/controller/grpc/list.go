package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) List(
	ctx context.Context,
	request *advertpb.ListRequest,
) (*advertpb.AdvertsResponse, error) {
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

	count, adverts, err := c.useCase.List(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to list adverts: %w", err)
	}

	response := advertsToResponse(adverts, count)

	return response, nil
}
