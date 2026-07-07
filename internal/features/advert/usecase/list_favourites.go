package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (uc *UseCase) ListFavourites(ctx context.Context, dto ListDTO) (int64, []advertentity.Advert, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return 0, nil, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if err := applyFilterScope(dto.UserID, dto.UserRole, &dto.Filter); err != nil {
		return 0, nil, fmt.Errorf("apply filter scope: %w", err)
	}

	favCount, err := uc.repo.CountFavourites(ctx, dto.UserID, dto.Filter)
	if err != nil {
		return 0, nil, fmt.Errorf("count favourites: %w", err)
	}

	favs, err := uc.repo.ListFavourites(ctx, dto.UserID, dto.Limit, dto.Offset, dto.Filter)
	if err != nil {
		return 0, nil, fmt.Errorf("list favourites: %w", err)
	}

	ids := make([]int64, len(favs))
	for i, a := range favs {
		ids[i] = a.ID
	}

	imagesByAdvertID, err := uc.repo.ListImagesByAdvertIDs(ctx, ids)
	if err != nil {
		return 0, nil, fmt.Errorf("list images by advert id: %w", err)
	}

	for i := range favs {
		advert := &favs[i]
		advert.Images = imagesByAdvertID[advert.ID]
	}

	return favCount, favs, nil
}
