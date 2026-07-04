package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) Archive(
	ctx context.Context,
	request *advertpb.ArchiveRequest,
) (*advertpb.AdvertResponse, error) {
	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.ArchiveDTO{
		UserID:   userID,
		UserRole: userRole,
		AdvertID: request.Id,
	}

	advert, err := c.useCase.Archive(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to archive advert with id = %d: %w", request.Id, err)
	}

	response := advertToResponse(advert)

	return response, nil
}
