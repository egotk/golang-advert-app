package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) ReissueRefreshToken(
	ctx context.Context,
	userID int,
	oldHash string,
	newToken userentity.RefreshToken,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", corepostgres.MapError(err))
	}
	defer tx.Rollback(ctx)

	deleteQuery := `
	DELETE FROM advertapp.refresh_tokens

	WHERE token_hash = $1 AND user_id = $2;
	`

	cmdTag, err := tx.Exec(ctx, deleteQuery, oldHash, userID)
	if err != nil {
		return fmt.Errorf("DELETE: %w", corepostgres.MapError(err))
	}
	if cmdTag.RowsAffected() == 0 {
		return userentity.ErrTokenNotFound
	}

	createQuery := `
	INSERT INTO advertapp.refresh_tokens (
		user_id,
		token_hash,
		issued_at,
		expires_at
	)

	VALUES ($1, $2, $3, $4);
	`

	if _, err := tx.Exec(
		ctx,
		createQuery,
		newToken.UserID,
		newToken.Hash,
		newToken.IssuedAt,
		newToken.ExpiresAt,
	); err != nil {
		return fmt.Errorf("INSERT: %w", corepostgres.MapError(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", corepostgres.MapError(err))
	}

	return nil
}
