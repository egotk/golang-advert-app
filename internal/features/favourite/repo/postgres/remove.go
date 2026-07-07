package favpostgres

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) Remove(
	ctx context.Context,
	advertID int64,
	userID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", corepostgres.MapError(err))
	}
	defer tx.Rollback(ctx)

	removeQuery := `
	DELETE FROM advertapp.favourites
	WHERE advert_id = $1 AND user_id = $2;
	`

	cmd, err := tx.Exec(ctx, removeQuery, advertID, userID)
	if err != nil {
		return corepostgres.MapError(err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("advert with id = %d not found in users favourites: %w", advertID, coreerrors.ErrNotFound)
	}

	decrementQuery := `
	UPDATE advertapp.adverts
	SET fav_count = fav_count - 1
	WHERE id = $1;
	`

	if _, err := tx.Exec(ctx, decrementQuery, advertID); err != nil {
		return corepostgres.MapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", corepostgres.MapError(err))
	}

	return nil

}
