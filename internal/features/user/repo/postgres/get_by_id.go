package userpostgres

import (
	"context"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) GetByID(
	ctx context.Context,
	id int,
) (userentity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, full_name, phone_number, role, locked_until, created_at, updated_at, image_path
	FROM advertapp.users
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var user userentity.User

	err := row.Scan(
		&user.ID,
		&user.Version,
		&user.Email,
		&user.FullName,
		&user.PhoneNumber,
		&user.Role,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ImagePath,
	)
	if err != nil {
		return userentity.User{}, corepostgres.MapError(err)
	}

	return user, nil
}
