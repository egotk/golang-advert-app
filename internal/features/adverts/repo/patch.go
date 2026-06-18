package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func (r *Repo) Patch(
	ctx context.Context,
	advert *advertentity.Advert,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE advertapp.adverts
	
	SET title = $1, description = $2, price = $3, category_id = $4, 
		status = $5, version = version + 1, updated_at = now()

	WHERE id = $6 AND version = $7

	RETURNING version, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.CategoryID,
		advert.Status,
		advert.ID,
		advert.Version,
	)
	err := row.Scan(
		&advert.Version,
		&advert.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return nil
}
