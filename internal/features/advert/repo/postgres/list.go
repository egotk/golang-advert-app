package advertpostgres

import (
	"context"
	"fmt"
	"strings"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) List(
	ctx context.Context,
	limit *int64,
	offset *int64,
	filter advertentity.Filter,
) ([]advertentity.Advert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id, version, user_id, title, description, price,
		category_id, status, views_count, created_at, updated_at

	FROM advertapp.adverts
	`)

	conditions, args := buildFilterConditions(filter)

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(orderBy(filter.Sort, filter.Order))

	args = append(args, limit)
	fmt.Fprintf(&queryBuilder, " LIMIT $%d", len(args))

	args = append(args, offset)
	fmt.Fprintf(&queryBuilder, " OFFSET $%d;", len(args))

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select: %w", corepostgres.MapError(err))
	}
	defer rows.Close()

	var adverts []advertentity.Advert
	for rows.Next() {
		var advert advertentity.Advert

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		adverts = append(adverts, advert)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", corepostgres.MapError(err))
	}

	return adverts, nil
}

func orderBy(
	sort *advertentity.Sort,
	order *advertentity.Order,
) string {
	if sort == nil {
		return " ORDER BY id ASC"
	}

	var column string
	switch *sort {
	case advertentity.SortByPrice:
		column = "price"
	case advertentity.SortByViews:
		column = "views_count"
	case advertentity.SortByDate:
		column = "created_at"
	default:
		return " ORDER BY id ASC"
	}

	direction := "ASC"
	if order != nil && *order == advertentity.OrderDesc {
		direction = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s, id ASC", column, direction)
}
