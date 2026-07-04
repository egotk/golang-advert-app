package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) Count(
	ctx context.Context,
	request *advertpb.CountRequest,
) (*advertpb.CountResponse, error) {
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

	count, err := c.useCase.Count(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to count adverts: %w", err)
	}

	response := &advertpb.CountResponse{Count: count}

	return response, nil
}
