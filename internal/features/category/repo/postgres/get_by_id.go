package categorypostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
)

func (r *Repo) GetByID(
	ctx context.Context,
	id int,
) (categoryentity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, parent_id, name

	FROM advertapp.advert_categories

	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var category categoryentity.Category
	err := row.Scan(
		&category.ID,
		&category.ParentID,
		&category.Name,
	)
	if err != nil {
		return categoryentity.Category{}, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return category, nil
}
