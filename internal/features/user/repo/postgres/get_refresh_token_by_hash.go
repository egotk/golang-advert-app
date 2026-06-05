package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) GetRefreshTokenByHash(
	ctx context.Context,
	hash string,
) (userentity.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT token_hash, user_id, issued_at, expires_at

	FROM advertapp.refresh_tokens

	WHERE token_hash = $1;
	`

	row := r.pool.QueryRow(ctx, query, hash)

	var token userentity.RefreshToken

	err := row.Scan(
		&token.Hash,
		&token.UserID,
		&token.IssuedAt,
		&token.ExpiresAt,
	)
	if err != nil {
		return userentity.RefreshToken{}, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return token, nil
}
