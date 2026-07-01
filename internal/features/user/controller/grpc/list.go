package usergrpc

import (
	"context"
	"fmt"

	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
)

func (c *Controller) List(
	ctx context.Context,
	r *userpb.ListRequest,
) (*userpb.UsersResponse, error) {
	users, err := c.useCase.ListUsers(ctx, r.Limit, r.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	response := usersToResponse(users)

	return response, nil
}
