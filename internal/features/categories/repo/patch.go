package categorypostgres

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/categories/entity"
)

func (r *Repo) Patch(
	ctx context.Context,
	category categoryentity.Category,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE advertapp.advert_categories

	SET parent_id = $1, name = $2

	WHERE id = $3
	`

	tag, err := r.pool.Exec(ctx, query, category.ParentID, category.Name, category.ID)
	if err != nil {
		return fmt.Errorf("update: %w", corepostgres.MapError(err))
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found: %w", coreerrors.ErrNotFound)
	}

	return nil
}
