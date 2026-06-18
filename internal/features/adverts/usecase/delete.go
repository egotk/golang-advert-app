package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/roles"
)

func (uc *UseCase) Delete(
	ctx context.Context,
	dto DeleteDTO,
) error {
	advert, err := uc.repo.GetByID(ctx, dto.AdvertID)
	if err != nil {
		return fmt.Errorf("get advert: %w", err)
	}

	if dto.UserID != advert.UserID && dto.UserRole != roles.Admin {
		return fmt.Errorf(
			"user cant delete others adverts: %w",
			coreerrors.ErrForbidden,
		)
	}

	if err := uc.repo.DeleteByID(ctx, dto.AdvertID); err != nil {
		return fmt.Errorf("delete advert: %w", err)
	}

	return nil
}
