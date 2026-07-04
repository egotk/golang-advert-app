package categoryusecase

import (
	"context"
	"fmt"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

func (uc *UseCase) Create(
	ctx context.Context,
	dto CreateDTO,
) (categoryentity.Category, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return categoryentity.Category{}, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if strings.TrimSpace(dto.Name) == "" {
		return categoryentity.Category{}, fmt.Errorf(
			"'Name' must not be empty: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	category := categoryentity.NewInitial(dto.ParentID, dto.Name)

	if err := uc.repo.Create(ctx, &category); err != nil {
		return categoryentity.Category{}, fmt.Errorf("create category: %w", err)
	}

	return category, nil
}
