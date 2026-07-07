package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) ListImagesByAdvertIDs(ctx context.Context, ids []int64) (map[int64][]advertentity.AdvertImage, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT advert_id, id, name, position, path, created_at

	FROM advertapp.advert_images

	WHERE advert_id = ANY($1)

	ORDER BY advert_id, position;
	`

	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("select: %w", corepostgres.MapError(err))
	}
	defer rows.Close()

	imagesByAdvertID := make(map[int64][]advertentity.AdvertImage)
	for rows.Next() {
		var advertID int64
		var image advertentity.AdvertImage

		err := rows.Scan(
			&advertID,
			&image.ID,
			&image.Name,
			&image.Position,
			&image.Path,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", corepostgres.MapError(err))
		}

		imagesByAdvertID[advertID] = append(imagesByAdvertID[advertID], image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", corepostgres.MapError(err))
	}

	return imagesByAdvertID, nil
}
