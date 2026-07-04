package advertgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) DeleteImage(
	ctx context.Context,
	request *advertpb.DeleteImageRequest,
) (*emptypb.Empty, error) {
	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dto := advertusecase.DeleteImageDTO{
		UserID:   userID,
		UserRole: userRole,
		ImageID:  request.Id,
	}

	if err := c.useCase.DeleteImage(ctx, dto); err != nil {
		return nil, fmt.Errorf("failed to delete image: %w", err)
	}

	return &emptypb.Empty{}, nil
}
