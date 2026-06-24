package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) Create(
	ctx context.Context,
	advert *advertentity.Advert,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", corepostgres.MapError(err))
	}
	defer tx.Rollback(ctx)

	advertQuery := `
	INSERT INTO advertapp.adverts (version, user_id, title, description, 
		price, category_id, status, views_count, created_at, updated_at)

	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)

	RETURNING id;
	`

	advertRow := tx.QueryRow(
		ctx,
		advertQuery,
		advert.Version,
		advert.UserID,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.CategoryID,
		advert.Status,
		advert.ViewsCount,
		advert.CreatedAt,
		advert.UpdatedAt,
	)
	if err := advertRow.Scan(&advert.ID); err != nil {
		return fmt.Errorf("scan advert: %w", corepostgres.MapError(err))
	}

	imageQuery := `
	INSERT INTO advertapp.advert_images (advert_id, name, position, path, created_at)

	VALUES ($1, $2, $3, $4, $5)

	RETURNING id, position;
	`

	for i := range advert.Images {
		image := &advert.Images[i]

		imageRow := tx.QueryRow(
			ctx,
			imageQuery,
			advert.ID,
			image.Name,
			i,
			image.Path,
			image.CreatedAt,
		)
		if err := imageRow.Scan(&image.ID, &image.Position); err != nil {
			return fmt.Errorf("scan image: %w", corepostgres.MapError(err))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", corepostgres.MapError(err))
	}

	return nil
}
