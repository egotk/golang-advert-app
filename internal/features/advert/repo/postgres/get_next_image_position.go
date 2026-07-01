package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	"github.com/jackc/pgx/v5"
)

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func getNextImagePosition(ctx context.Context, rower queryRower, advertID int64) (int64, error) {
	const posQuery = `
	SELECT COALESCE(MAX(position), -1) + 1

	FROM advertapp.advert_images

	WHERE advert_id = $1;
	`

	var nextPosition int64
	row := rower.QueryRow(ctx, posQuery, advertID)
	if err := row.Scan(&nextPosition); err != nil {
		return 0, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return nextPosition, nil
}
