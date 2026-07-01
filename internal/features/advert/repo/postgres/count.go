package advertpostgres

import (
	"context"
	"fmt"
	"strings"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) Count(
	ctx context.Context,
	filter advertentity.Filter,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT COUNT(*)

	FROM advertapp.adverts
	`)

	conditions, args := buildFilterConditions(filter)

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}
	queryBuilder.WriteString(";")

	var count int64
	row := r.pool.QueryRow(ctx, queryBuilder.String(), args...)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("select count: %w", corepostgres.MapError(err))
	}

	return count, nil
}
