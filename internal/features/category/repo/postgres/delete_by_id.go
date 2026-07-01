package categorypostgres

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) DeleteByID(
	ctx context.Context,
	id int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE 

	FROM advertapp.advert_categories

	WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete: %w", corepostgres.MapError(err))
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found: %w", coreerrors.ErrNotFound)
	}

	return nil
}
