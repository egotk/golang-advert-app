package categorygrpc

import (
	"context"
	"fmt"

	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Delete(
	ctx context.Context,
	request *categorypb.DeleteRequest,
) (*emptypb.Empty, error) {
	if err := c.useCase.Delete(ctx, request.Id); err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	return &emptypb.Empty{}, nil
}
