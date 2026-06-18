package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (uc *UseCase) GetByID(
	ctx context.Context,
	dto GetByIDDTO,
) (advertentity.Advert, error) {
	if dto.AdvertID <= 0 {
		return advertentity.Advert{}, fmt.Errorf("'ID' must be positive: %w", coreerrors.ErrInvalidArgument)
	}

	advert, err := uc.repo.GetByID(ctx, dto.AdvertID)
	if err != nil {
		return advertentity.Advert{}, fmt.Errorf("get advert by ID: %w", err)
	}

	// запретить всем кроме владельца и админов просмотр объявлений со статусом != active
	if advert.Status != advertentity.StatusActive && dto.UserID != advert.UserID {
		if dto.UserRole != roles.Admin {
			return advertentity.Advert{}, fmt.Errorf("invalid role: %w", coreerrors.ErrForbidden)
		}
	}

	if err := uc.repo.IncrementViewsCount(ctx, dto.AdvertID); err != nil {
		return advertentity.Advert{}, fmt.Errorf("increment views count: %w", err)
	}

	return advert, nil
}
