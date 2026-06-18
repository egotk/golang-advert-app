package advertpostgres

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) IncrementViewsCount(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE advertapp.adverts

	SET views_count = views_count + 1

	WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec: %w", corepostgres.MapError(err))
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("advert not found: %w", coreerrors.ErrNotFound)
	}

	return nil
}
