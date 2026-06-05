package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) CreateRefreshToken(
	ctx context.Context,
	token userentity.RefreshToken,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO advertapp.refresh_tokens (
		user_id,
		token_hash,
		issued_at,
		expires_at
	)

	VALUES ($1, $2, $3, $4);
	`

	if _, err := r.pool.Exec(
		ctx,
		query,
		token.UserID,
		token.Hash,
		token.IssuedAt,
		token.ExpiresAt,
	); err != nil {
		return fmt.Errorf("INSERT: %w", corepostgres.MapError(err))
	}

	return nil
}
