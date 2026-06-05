package userpostgres

import (
	"context"
	"fmt"
	"time"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) IncrementFailedLoginCount(
	ctx context.Context,
	id int,
	version int,
) (*time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE advertapp.users

	SET 
		version = version + 1,
		failed_login_count = failed_login_count + 1,
		locked_until = CASE
			WHEN failed_login_count + 1 >= 5
			THEN now() + INTERVAL '5 minutes'
			ELSE locked_until
		END

	WHERE id = $1 AND version = $2

	RETURNING locked_until;
	`

	row := r.pool.QueryRow(ctx, query, id, version)

	var lockedUntil *time.Time

	err := row.Scan(
		&lockedUntil,
	)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return lockedUntil, nil
}
