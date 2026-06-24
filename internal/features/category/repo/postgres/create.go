package categorypostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

func (r *Repo) Create(
	ctx context.Context,
	category *categoryentity.Category,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO advertapp.advert_categories (parent_id, name)

	VALUES ($1, $2)

	RETURNING id;
	`

	row := r.pool.QueryRow(ctx, query, category.ParentID, category.Name)
	if err := row.Scan(&category.ID); err != nil {
		return fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return nil
}
