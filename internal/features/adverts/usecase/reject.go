package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (uc *UseCase) Reject(
	ctx context.Context,
	id int,
) (advertentity.Advert, error) {
	if id <= 0 {
		return advertentity.Advert{}, fmt.Errorf("'ID' must be positive: %w", coreerrors.ErrInvalidArgument)
	}

	advert, err := uc.repo.ChangeStatus(
		ctx,
		id,
		advertentity.StatusInitial,
		advertentity.StatusRejected,
	)
	if err != nil {
		return advertentity.Advert{}, fmt.Errorf("reject advert: %w", err)
	}

	return advert, nil
}
