package advertpostgres

import (
	"context"
	"fmt"
	"strings"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) ListFavourites(
	ctx context.Context,
	userID int64,
	limit *int64,
	offset *int64,
	filter advertentity.Filter,
) ([]advertentity.Advert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id, version, a.user_id, title, description, price, category_id, 
		status, views_count, a.created_at, updated_at, fav_count
	FROM advertapp.adverts a
	JOIN advertapp.favourites f ON f.advert_id = a.id AND f.user_id = $1
	`)

	tableName := "a"
	conditions, args := buildAdvertFilterConditions(filter, &tableName, []any{userID})

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(orderBy(filter.Sort, &tableName, filter.Order))

	args = append(args, limit)
	fmt.Fprintf(&queryBuilder, " LIMIT $%d", len(args))

	args = append(args, offset)
	fmt.Fprintf(&queryBuilder, " OFFSET $%d;", len(args))

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select: %w", corepostgres.MapError(err))
	}
	defer rows.Close()

	var favs []advertentity.Advert
	for rows.Next() {
		var fav advertentity.Advert

		err := rows.Scan(
			&fav.ID,
			&fav.Version,
			&fav.UserID,
			&fav.Title,
			&fav.Description,
			&fav.Price,
			&fav.CategoryID,
			&fav.Status,
			&fav.ViewsCount,
			&fav.CreatedAt,
			&fav.UpdatedAt,
			&fav.FavouriteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		favs = append(favs, fav)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", corepostgres.MapError(err))
	}

	return favs, nil
}
