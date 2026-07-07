package favusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
)

func (uc *UseCase) Remove(ctx context.Context, dto RemoveDTO) error {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if err := uc.repo.Remove(ctx, dto.AdvertID, dto.UserID); err != nil {
		return fmt.Errorf("remove from favourites: %w", err)
	}

	return nil
}
