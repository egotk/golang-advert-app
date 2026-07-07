package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) GetByID(ctx context.Context, id int64) (advertentity.Advert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT 
		id,
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

	FROM advertapp.adverts
	
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var advert advertentity.Advert

	err := row.Scan(
		&advert.ID,
		&advert.Version,
		&advert.UserID,
		&advert.Title,
		&advert.Description,
		&advert.Price,
		&advert.CategoryID,
		&advert.Status,
		&advert.ViewsCount,
		&advert.CreatedAt,
		&advert.UpdatedAt,
	)
	if err != nil {
		return advertentity.Advert{}, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return advert, nil
}
