package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (uc *UseCase) List(ctx context.Context, dto ListDTO) (int64, []advertentity.Advert, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return 0, nil, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if err := applyFilterScope(dto.UserID, dto.UserRole, &dto.Filter); err != nil {
		return 0, nil, fmt.Errorf("apply filter scope: %w", err)
	}

	count, err := uc.repo.Count(ctx, dto.Filter)
	if err != nil {
		return 0, nil, fmt.Errorf("count adverts: %w", err)
	}

	adverts, err := uc.repo.List(ctx, dto.Limit, dto.Offset, dto.Filter)
	if err != nil {
		return 0, nil, fmt.Errorf("get adverts from repo: %w", err)
	}

	ids := make([]int64, len(adverts))
	for i, a := range adverts {
		ids[i] = a.ID
	}

	imagesByAdvertID, err := uc.repo.ListImagesByAdvertIDs(ctx, ids)
	if err != nil {
		return 0, nil, fmt.Errorf("list images by advert id: %w", err)
	}

	for i := range adverts {
		advert := &adverts[i]
		advert.Images = imagesByAdvertID[advert.ID]
	}

	return count, adverts, nil
}
