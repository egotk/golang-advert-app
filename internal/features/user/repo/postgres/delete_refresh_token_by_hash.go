package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
)

func (r *Repo) DeleteRefreshTokenByHash(
	ctx context.Context,
	userId int64,
	hash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM advertapp.refresh_tokens
	
	WHERE token_hash = $1 AND user_id = $2;
	`

	if _, err := r.pool.Exec(ctx, query, hash, userId); err != nil {
		return fmt.Errorf("DELETE: %w", corepostgres.MapError(err))
	}

	return nil
}
