package categoryusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

func (uc *UseCase) Patch(
	ctx context.Context,
	dto PatchDTO,
) (categoryentity.Category, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return categoryentity.Category{}, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if dto.ParentID.Value != nil && *dto.ParentID.Value < 1 {
		return categoryentity.Category{}, fmt.Errorf("'ParentID' must be positive: %w", coreerrors.ErrInvalidArgument)
	}

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
