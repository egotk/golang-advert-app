package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
)

func (uc *UseCase) AddToFavourites(ctx context.Context, dto AddToFavouritesDTO) error {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	advert, err := uc.repo.GetByID(ctx, dto.AdvertID)
	if err != nil {
		return fmt.Errorf("get advert: %w", err)
	}

	if dto.UserID != advert.UserID && dto.UserRole != roles.Admin {
		if !advert.IsPublic() {
			return coreerrors.ErrForbidden
		}
	}

	if err := uc.favRepo.Add(ctx, dto.AdvertID, dto.UserID); err != nil {
		return fmt.Errorf("add to favourites: %w", err)
	}

	return nil
}
