package usergrpc

import (
	"context"
	"fmt"

	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

func (c *Controller) Login(
	ctx context.Context,
	request *userpb.LoginRequest,
) (*userpb.LoginResponse, error) {
	dto := userusecase.LoginDTO{
		Email:    request.Email,
		Password: request.Password,
	}

	result, err := c.useCase.Login(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	response := &userpb.LoginResponse{
		UserId:       result.UserID,
		AccessToken:  result.Tokens.Access,
		RefreshToken: result.Tokens.Refresh,
	}

	return response, nil
}
