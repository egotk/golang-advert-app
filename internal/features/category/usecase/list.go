package categoryusecase

import (
	"context"
	"fmt"

	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

func (uc *UseCase) List(ctx context.Context) ([]categoryentity.Category, error) {
	categories, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get categories from repo: %w", err)
	}

	return categories, nil
}
