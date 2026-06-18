package categorypostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/categories/entity"
)

func (r *Repo) List(ctx context.Context) ([]categoryentity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, parent_id, name
	
	FROM advertapp.advert_categories

	ORDER BY name ASC, id ASC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select: %w", corepostgres.MapError(err))
	}
	defer rows.Close()

	var categories []categoryentity.Category
	for rows.Next() {
		var category categoryentity.Category

		err := rows.Scan(
			&category.ID,
			&category.ParentID,
			&category.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", corepostgres.MapError(err))
	}

	return categories, nil
}
