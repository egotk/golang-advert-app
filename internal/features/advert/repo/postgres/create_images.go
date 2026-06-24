package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) CreateImages(
	ctx context.Context,
	advertID int,
	images []advertentity.AdvertImage,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", corepostgres.MapError(err))
	}
	defer tx.Rollback(ctx)

	nextPosition, err := getNextImagePosition(ctx, tx, advertID)
	if err != nil {
		return fmt.Errorf("get next image position: %w", err)
	}

	const query = `
	INSERT INTO advertapp.advert_images (advert_id, 
		name, position, path, created_at)

	VALUES ($1, $2, $3, $4, $5)

	RETURNING id, position;
	`

	for i := range images {
		img := &images[i]

		row := tx.QueryRow(
			ctx,
			query,
			advertID,
			img.Name,
			nextPosition,
			img.Path,
			img.CreatedAt,
		)
		if err := row.Scan(&img.ID, &img.Position); err != nil {
			return fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		nextPosition++
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", corepostgres.MapError(err))
	}

	return nil
}
