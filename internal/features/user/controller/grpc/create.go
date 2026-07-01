package usergrpc

import (
	"context"
	"fmt"

	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

func (c *Controller) Create(
	ctx context.Context,
	r *userpb.CreateRequest,
) (*userpb.UserResponse, error) {
	dto := userusecase.CreateDTO{
		Email:       r.Email,
		FullName:    r.FullName,
		PhoneNumber: r.PhoneNumber,
		Password:    r.Password,
	}

	user, err := c.useCase.CreateUser(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := userToResponse(user)

	return response, nil
}
