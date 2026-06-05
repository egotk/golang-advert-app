package userusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (uc *UseCase) ListUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]userentity.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	users, err := uc.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users from repo: %w", err)
	}

	return users, nil
}
