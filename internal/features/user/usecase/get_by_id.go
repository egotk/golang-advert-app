package userusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (uc *UseCase) GetByID(
	ctx context.Context,
	id int,
) (userentity.User, error) {
	if id < 0 {
		return userentity.User{}, fmt.Errorf(
			"id must be positive: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return userentity.User{}, fmt.Errorf(
			"get user with id = %d from repo: %w",
			id,
			err,
		)
	}

	return user, nil
}
