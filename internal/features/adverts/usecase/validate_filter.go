package advertusecase

import (
	"fmt"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func validateFilter(filter advertentity.Filter) error {
	if filter.Title != nil && strings.TrimSpace(*filter.Title) == "" {
		return fmt.Errorf(
			"'Title' must not be empty: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if filter.Description != nil && strings.TrimSpace(*filter.Description) == "" {
		return fmt.Errorf(
			"'Description' must not be empty: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if filter.MinPrice != nil && *filter.MinPrice < 0 {
		return fmt.Errorf(
			"'MinPrice' must be non negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if filter.MaxPrice != nil && *filter.MaxPrice < 0 {
		return fmt.Errorf(
			"'MaxPrice' must be non negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
		return fmt.Errorf(
			"'MinPrice' must be lower than max price: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if filter.CategoryID != nil && *filter.CategoryID <= 0 {
		return fmt.Errorf(
			"'CategoryID' must be positive: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if filter.FromDate != nil && filter.ToDate != nil && filter.FromDate.After(*filter.ToDate) {
		return fmt.Errorf(
			"'FromDate' must be earlier than 'ToDate': %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}
