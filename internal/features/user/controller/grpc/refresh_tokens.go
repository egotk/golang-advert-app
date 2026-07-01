package usergrpc

import (
	"context"
	"fmt"

	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

func (c *Controller) RefreshTokens(
	ctx context.Context,
	request *userpb.RefreshTokensRequest,
) (*userpb.RefreshTokensResponse, error) {
	dto := userusecase.RefreshTokensDTO{
		RefreshToken: request.RefreshToken,
	}

	tokens, err := c.useCase.RefreshTokens(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}

	response := &userpb.RefreshTokensResponse{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}

	return response, nil
}
