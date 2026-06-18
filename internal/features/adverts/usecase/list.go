package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (uc *UseCase) List(
	ctx context.Context,
	dto ListDTO,
) (int, []advertentity.Advert, error) {
	if dto.Limit != nil && *dto.Limit < 0 {
		return 0, nil, fmt.Errorf(
			"'Limit' must be non negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if dto.Offset != nil && *dto.Offset < 0 {
		return 0, nil, fmt.Errorf(
			"'Offset' must be non negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if err := validateFilter(dto.Filter); err != nil {
		return 0, nil, fmt.Errorf("validate filter: %w", err)
	}

	if err := applyFilterScope(dto.UserRole, &dto.Filter); err != nil {
		return 0, nil, fmt.Errorf("apply filter scope: %w", err)
	}

	adverts, err := uc.repo.List(ctx, dto.Limit, dto.Offset, dto.Filter)
	if err != nil {
		return 0, nil, fmt.Errorf("get adverts from repo: %w", err)
	}

	count, err := uc.repo.Count(ctx, dto.Filter)
	if err != nil {
		return 0, nil, fmt.Errorf("count adverts: %w", err)
	}

	return count, adverts, nil
}
