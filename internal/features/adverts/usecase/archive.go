package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (uc *UseCase) Archive(
	ctx context.Context,
	dto ArchiveDTO,
) (advertentity.Advert, error) {
	if dto.AdvertID <= 0 {
		return advertentity.Advert{}, fmt.Errorf("'ID' must be positive: %w", coreerrors.ErrInvalidArgument)
	}

	advert, err := uc.repo.GetByID(ctx, dto.AdvertID)
	if err != nil {
		return advertentity.Advert{}, fmt.Errorf("get advert by ID: %w", err)
	}

	if advert.UserID != dto.UserID && dto.UserRole != roles.Admin {
		return advertentity.Advert{}, fmt.Errorf("can't archive other's advert: %w", coreerrors.ErrForbidden)
	}

	result, err := uc.repo.ChangeStatus(
		ctx,
		dto.AdvertID,
		advertentity.StatusActive,
		advertentity.StatusArchived,
	)
	if err != nil {
		return advertentity.Advert{}, fmt.Errorf("archive advert: %w", err)
	}

	return result, nil
}
