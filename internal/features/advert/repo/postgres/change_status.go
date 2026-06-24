package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) ChangeStatus(
	ctx context.Context,
	id int,
	oldStatus advertentity.Status,
	newStatus advertentity.Status,
) (advertentity.Advert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE advertapp.adverts

	SET status = $1, version = version + 1, updated_at = now()

	WHERE id = $2 AND status = $3

	RETURNING
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
		updated_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		string(newStatus),
		id,
		string(oldStatus),
	)

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
