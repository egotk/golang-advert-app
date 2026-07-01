package usergrpc

import (
	"context"
	"fmt"

	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

func (c *Controller) GetByID(
	ctx context.Context,
	request *userpb.GetByIDRequest,
) (*userpb.UserResponse, error) {
	user, err := c.useCase.GetUserByID(ctx, request.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	response := userToResponse(user)

	return response, nil
}
