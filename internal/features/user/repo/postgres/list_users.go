package userpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (r *Repo) ListUsers(
	ctx context.Context,
	limit *int64,
	offset *int64,
) ([]userentity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, full_name, phone_number, role, locked_until, created_at, updated_at, image_path

	FROM advertapp.users

	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SELECT: %w", corepostgres.MapError(err))
	}
	defer rows.Close()

	var users []userentity.User
	for rows.Next() {
		var user userentity.User

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", corepostgres.MapError(err))
	}

	return users, nil
}
