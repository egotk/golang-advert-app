package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Delete(
	ctx context.Context,
	request *advertpb.DeleteRequest,
) (*emptypb.Empty, error) {
	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.DeleteDTO{
		UserID:   userID,
		UserRole: userRole,
		AdvertID: request.Id,
	}

	if err := c.useCase.Delete(ctx, dto); err != nil {
		return nil, fmt.Errorf("failed to delete advert: %w", err)
	}

	return &emptypb.Empty{}, nil
}
