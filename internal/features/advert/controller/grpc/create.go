package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) Create(
	ctx context.Context,
	request *advertpb.CreateRequest,
) (*advertpb.AdvertResponse, error) {
	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.CreateDTO{
		UserID:      userID,
		Title:       request.Title,
		Description: request.Description,
		Price:       request.Price,
		CategoryID:  request.CategoryId,
	}

	advert, err := c.useCase.Create(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to create advert: %w", err)
	}

	response := advertToResponse(advert)

	return response, nil
}
