package favpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) Add(
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

	addQuery := `
	INSERT INTO advertapp.favourites(advert_id, user_id, created_at)
	VALUES ($1, $2, now())
	ON CONFLICT (advert_id, user_id) DO NOTHING;
	`

	tag, err := tx.Exec(ctx, addQuery, advertID, userID)
	if err != nil {
		return corepostgres.MapError(err)
	}

	if tag.RowsAffected() != 0 {
		incrementQuery := `
		UPDATE advertapp.adverts
		SET fav_count = fav_count + 1
		WHERE id = $1;
	`

		if _, err := tx.Exec(ctx, incrementQuery, advertID); err != nil {
			return corepostgres.MapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", corepostgres.MapError(err))
	}

	return nil
}
