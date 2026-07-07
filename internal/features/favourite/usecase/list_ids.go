package favusecase

import (
	"context"
	"fmt"
)

func (uc *UseCase) ListIDs(ctx context.Context, userID int64) ([]int64, error) {
	ids, err := uc.repo.ListIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list favourite ids: %w", err)
	}

	return ids, nil
}
