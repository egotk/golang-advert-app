package advertpostgres

import (
	"context"
	"fmt"

	corepostgres "github.com/egotk/golang-advert-app/internal/core/postgres"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (r *Repo) GetImageByID(
	ctx context.Context,
	imageID int64,
) (int64, advertentity.AdvertImage, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT advert_id, id, name, position, path, created_at

	FROM advertapp.advert_images

	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, imageID)

	var advertID int64
	var image advertentity.AdvertImage
	err := row.Scan(
		&advertID,
		&image.ID,
		&image.Name,
		&image.Position,
		&image.Path,
		&image.CreatedAt,
	)
	if err != nil {
		return 0, advertentity.AdvertImage{}, fmt.Errorf("scan: %w", corepostgres.MapError(err))
	}

	return advertID, image, nil
}
