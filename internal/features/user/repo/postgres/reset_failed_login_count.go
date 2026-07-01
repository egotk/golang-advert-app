package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) ResetFailedLoginCount(
	ctx context.Context,
	id int64,
	version int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE advertapp.users

	SET
		version = version + 1,
		failed_login_count = 0,
		locked_until = NULL
	
	WHERE id = $1 AND version = $2;
	`

	tag, err := r.pool.Exec(ctx, query, id, version)
	if err != nil {
		return fmt.Errorf("UPDATE: %w", corepostgres.MapError(err))
	}
	if tag.RowsAffected() == 0 {
		return userentity.ErrUserVersionConflict
	}

	return nil
}
