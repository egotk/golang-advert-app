package advertusecase

import (
	"context"
	"fmt"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (uc *UseCase) Patch(
	ctx context.Context,
	dto PatchDTO,
) (advertentity.Advert, error) {
	if dto.Title == nil && dto.Description == nil && dto.Price == nil && dto.CategoryID == nil {
		return advertentity.Advert{}, fmt.Errorf(
			"empty patch request: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if dto.Title != nil && strings.TrimSpace(*dto.Title) == "" {
		return advertentity.Advert{}, fmt.Errorf(
			"'Title' must not be empty: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if dto.Description != nil && strings.TrimSpace(*dto.Description) == "" {
		return advertentity.Advert{}, fmt.Errorf(
			"'Description' must not be empty: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	advert, err := uc.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return advertentity.Advert{}, fmt.Errorf("get advert: %w", err)
	}

	if advert.UserID != dto.UserID && dto.UserRole == roles.User {
		return advertentity.Advert{}, fmt.Errorf(
			"user cant patch others adverts: %w",
			coreerrors.ErrForbidden,
		)
	}

	advert.Version = dto.Version

	if dto.Title != nil {
		advert.Title = *dto.Title
	}

	if dto.Description != nil {
		advert.Description = *dto.Description
	}

	if dto.Price != nil {
		advert.Price = *dto.Price
	}

	if dto.CategoryID != nil {
		advert.CategoryID = *dto.CategoryID
	}

	if dto.UserRole != roles.Admin {
		if advert.Status == advertentity.StatusBlocked {
			return advertentity.Advert{}, fmt.Errorf(
				"cant patch blocked advert: %w",
				coreerrors.ErrForbidden,
			)
		}

		if dto.Title != nil || dto.Description != nil {
			advert.Status = advertentity.StatusInitial
		}
	}

	if err := uc.repo.Patch(ctx, &advert); err != nil {
		return advertentity.Advert{}, fmt.Errorf("patch advert: %w", err)
	}

	return advert, nil
}
