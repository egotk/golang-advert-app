package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) GetUserByEmail(
	ctx context.Context,
	email string,
) (userentity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT 
		id,
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

	FROM advertapp.users
	
	WHERE email = $1;
	`

	row := r.pool.QueryRow(ctx, query, email)

	var user userentity.User
	err := row.Scan(
		&user.ID,
		&user.Version,
		&user.Email,
		&user.FullName,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Role,
		&user.FailedLoginCount,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ImagePath,
	)
	if err != nil {
		return userentity.User{}, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return user, nil
}
