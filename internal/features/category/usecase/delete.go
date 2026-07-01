package categoryusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
)

func (uc *UseCase) Delete(
	ctx context.Context,
	id int64,
) error {
	if id < 1 {
		return fmt.Errorf("'ID' must be positive: %w", coreerrors.ErrInvalidArgument)
	}

	if err := uc.repo.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
