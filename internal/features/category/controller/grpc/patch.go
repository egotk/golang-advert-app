package categorygrpc

import (
	"context"
	"fmt"

	"github.com/egotk/golang-advert-app/internal/core/nullable"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
)

type patchRequest struct {
	Id       int64
	ParentID nullable.Nullable[int64]
	Name     *string
}

func patchRequestfromGRPC(request *categorypb.PatchRequest) patchRequest {
	var parentID nullable.Nullable[int64]
	switch v := request.ParentIdOpt.(type) {
	case *categorypb.PatchRequest_ParentId:
		parentID = nullable.Nullable[int64]{
			Set:   true,
			Value: &v.ParentId,
		}
	case *categorypb.PatchRequest_ShouldClear:
		parentID = nullable.Nullable[int64]{
			Set:   true,
			Value: nil,
		}
	}

	return patchRequest{
		Id:       request.Id,
		ParentID: parentID,
		Name:     request.Name,
	}
}

func (r patchRequest) toDTO() categoryusecase.PatchDTO {
	return categoryusecase.PatchDTO{
		ID:       r.Id,
		ParentID: r.ParentID,
		Name:     r.Name,
	}
}

func (c *Controller) Patch(
	ctx context.Context,
	request *categorypb.PatchRequest,
) (*categorypb.CategoryResponse, error) {
	dto := patchRequestfromGRPC(request).toDTO()

	category, err := c.useCase.Patch(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to patch category: %w", err)
	}

	response := categoryToResponse(category)

	return response, nil
}
