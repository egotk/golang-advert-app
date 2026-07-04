package categorygrpc

import (
	"context"
	"fmt"

	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) List(
	ctx context.Context,
	_ *emptypb.Empty,
) (*categorypb.CategoriesResponse, error) {
	categories, err := c.useCase.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	response := categoriesToResponse(categories)

	return response, nil
}
