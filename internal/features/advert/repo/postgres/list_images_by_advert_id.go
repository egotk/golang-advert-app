package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) ListImagesByAdvertID(ctx context.Context, advertID int64) ([]advertentity.AdvertImage, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, name, position, path, created_at

	FROM advertapp.advert_images

	WHERE advert_id = $1

	ORDER BY position ASC;
	`

	rows, err := r.pool.Query(ctx, query, advertID)
	if err != nil {
		return nil, fmt.Errorf("select: %w", corepostgres.MapError(err))
	}
	defer rows.Close()

	var images []advertentity.AdvertImage
	for rows.Next() {
		var image advertentity.AdvertImage

		err := rows.Scan(
			&image.ID,
			&image.Name,
			&image.Position,
			&image.Path,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", corepostgres.MapError(err))
	}

	return images, nil
}
