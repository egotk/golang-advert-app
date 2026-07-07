package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
)

func (uc *UseCase) CountFavourites(ctx context.Context, dto CountDTO) (int64, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return 0, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if err := applyFilterScope(dto.UserID, dto.UserRole, &dto.Filter); err != nil {
		return 0, fmt.Errorf("apply filter scope: %w", err)
	}

	favCount, err := uc.repo.CountFavourites(ctx, dto.UserID, dto.Filter)
	if err != nil {
		return 0, fmt.Errorf("count favourites: %w", err)
	}

	return favCount, nil
}
