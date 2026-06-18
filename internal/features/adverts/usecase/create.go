package advertusecase

import (
	"context"
	"fmt"

	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (uc *UseCase) Create(
	ctx context.Context,
	dto CreateDTO,
) (advertentity.Advert, error) {
	advert := advertentity.NewInitial(
		dto.UserID,
		dto.Title,
		dto.Description,
		dto.Price,
		dto.CategoryID,
	)

	if err := uc.repo.Create(ctx, &advert); err != nil {
		return advertentity.Advert{}, fmt.Errorf("store advert in DB: %w", err)
	}

	return advert, nil
}
