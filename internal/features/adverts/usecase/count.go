package advertusecase

import (
	"context"
	"fmt"
)

func (uc *UseCase) Count(
	ctx context.Context,
	dto CountDTO,
) (int, error) {
	if err := validateFilter(dto.Filter); err != nil {
		return 0, fmt.Errorf("validate filter: %w", err)
	}

	if err := applyFilterScope(dto.UserRole, &dto.Filter); err != nil {
		return 0, fmt.Errorf("apply filter scope: %w", err)
	}

	count, err := uc.repo.Count(ctx, dto.Filter)
	if err != nil {
		return 0, fmt.Errorf("count adverts: %w", err)
	}

	return count, nil
}
