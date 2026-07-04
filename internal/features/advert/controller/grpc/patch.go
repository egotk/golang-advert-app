package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) Patch(
	ctx context.Context,
	request *advertpb.PatchRequest,
) (*advertpb.AdvertResponse, error) {
	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.PatchDTO{
		UserID:      userID,
		UserRole:    userRole,
		ID:          request.Id,
		Version:     request.Version,
		Title:       request.Title,
		Description: request.Description,
		Price:       request.Price,
		CategoryID:  request.CategoryId,
	}

	advert, err := c.useCase.Patch(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to patch advert: %w", err)
	}

	response := advertToResponse(advert)

	return response, nil
}
