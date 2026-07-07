package advertpostgres

import (
	"context"
	"fmt"
	"strings"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) CountFavourites(
	ctx context.Context,
	userID int64,
	filter advertentity.Filter,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT COUNT(*)
	FROM advertapp.adverts a
	JOIN advertapp.favourites f ON f.advert_id = a.id AND f.user_id = $1
	`)

	tableName := "a"
	conditions, args := buildAdvertFilterConditions(filter, &tableName, []any{userID})

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}
	queryBuilder.WriteString(";")

	var favCount int64
	row := r.pool.QueryRow(ctx, queryBuilder.String(), args...)
	if err := row.Scan(&favCount); err != nil {
		return 0, fmt.Errorf("select count: %w", corepostgres.MapError(err))
	}

	return favCount, nil
}
