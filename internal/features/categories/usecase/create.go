package categoryusecase

import (
	"context"
	"fmt"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/categories/entity"
)

func (uc *UseCase) Create(
	ctx context.Context,
	dto CreateDTO,
) (categoryentity.Category, error) {
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
