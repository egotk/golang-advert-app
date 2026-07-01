package userusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (uc *UseCase) GetUserByID(
	ctx context.Context,
	id int64,
) (userentity.User, error) {
	if id <= 0 {
		return userentity.User{}, fmt.Errorf(
			"id must be positive: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return userentity.User{}, fmt.Errorf(
			"get user with id = %d from repo: %w",
			id,
			err,
		)
	}

	return user, nil
}
