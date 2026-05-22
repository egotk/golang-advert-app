package userpostgres

import (
	"context"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) Create(
	ctx context.Context,
	user *userentity.User,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO advertapp.users (
		version,
		email,
		full_name,
		phone_number,
		password_hash,
		role,
		failed_login_count,
		locked_until,
		created_at,
		updated_at,
		image_path
	)
	
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, version;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Version,
		user.Email,
		user.FullName,
		user.PhoneNumber,
		user.PasswordHash,
		user.Role,
		user.FailedLoginCount,
		user.LockedUntil,
		user.CreatedAt,
		user.UpdatedAt,
		user.ImagePath,
	)

	err := row.Scan(
		&user.ID,
		&user.Version,
	)
	if err != nil {
		return corepostgres.MapError(err)
	}

	return nil
}
