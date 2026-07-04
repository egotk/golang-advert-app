package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) GetByID(
	ctx context.Context,
	request *advertpb.GetByIDRequest,
) (*advertpb.AdvertResponse, error) {
	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.GetByIDDTO{
		UserID:   userID,
		UserRole: userRole,
		AdvertID: request.Id,
	}

	advert, err := c.useCase.GetByID(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get request by id: %w", err)
	}

	response := advertToResponse(advert)

	return response, nil
}
