package categorygrpc

import (
	"context"
	"fmt"

	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
)

func (c *Controller) Create(
	ctx context.Context,
	request *categorypb.CreateRequest,
) (*categorypb.CategoryResponse, error) {
	dto := categoryusecase.CreateDTO{
		ParentID: request.ParentId,
		Name:     request.Name,
	}

	category, err := c.useCase.Create(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	response := categoryToResponse(category)

	return response, nil
}
