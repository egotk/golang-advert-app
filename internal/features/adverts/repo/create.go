package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (r *Repo) Create(
	ctx context.Context,
	advert *advertentity.Advert,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO advertapp.adverts (
		version,
		user_id,
		title,
		description,
		price,
		category_id,
		status,
		views_count,
		created_at,
		updated_at
	)

	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)

	RETURNING id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
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

	if err := row.Scan(&advert.ID); err != nil {
		return fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return nil
}
