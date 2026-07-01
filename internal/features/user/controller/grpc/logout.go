package usergrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Logout(
	ctx context.Context,
	request *userpb.LogoutRequest,
) (*emptypb.Empty, error) {
	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get 'UserID' from JWT: %w", err)
	}

	dto := userusecase.LogoutDTO{
		UserID:       userID,
		RefreshToken: request.RefreshToken,
	}

	if err := c.useCase.Logout(ctx, dto); err != nil {
		return nil, fmt.Errorf("failed to logout: %w", err)
	}

	return &emptypb.Empty{}, nil
}
