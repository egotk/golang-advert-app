package advertpostgres

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) DeleteImageByID(
	ctx context.Context,
	imageID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM advertapp.advert_images

	WHERE id = $1;
	`

	cmd, err := r.pool.Exec(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("delete: %w", corepostgres.MapError(err))
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf(
			"image with id=%d not found: %w",
			imageID,
			coreerrors.ErrNotFound,
		)
	}

	return nil
}
