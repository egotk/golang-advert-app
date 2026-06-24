package categoryusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

func (uc *UseCase) Patch(
	ctx context.Context,
	dto PatchDTO,
) (categoryentity.Category, error) {
	if !dto.ParentID.Set && dto.Name == nil {
		return categoryentity.Category{}, fmt.Errorf(
			"empty patch request: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	category, err := uc.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return categoryentity.Category{}, fmt.Errorf("get category from repo: %w", err)
	}

	if dto.ParentID.Set {
		if dto.ParentID.Value != nil && *dto.ParentID.Value == dto.ID {
			return categoryentity.Category{}, fmt.Errorf(
				"category cant be its own parent: %w",
				coreerrors.ErrInvalidArgument,
			)
		}

		category.ParentID = dto.ParentID.Value
	}

	if dto.Name != nil {
		category.Name = *dto.Name
	}

	if err := uc.repo.Patch(ctx, category); err != nil {
		return categoryentity.Category{}, fmt.Errorf("patch category: %w", err)
	}

	return category, nil
}
