package advertusecase

import (
	"context"
	"fmt"
	"io"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func (uc *UseCase) GetImageByID(
	ctx context.Context,
	dto GetImageDTO,
) (io.ReadCloser, advertentity.AdvertImage, error) {
	advertID, advertImage, err := uc.repo.GetImageByID(ctx, dto.ImageID)
	if err != nil {
		return nil, advertentity.AdvertImage{},
			fmt.Errorf("get image from repo: %w", err)
	}

	advert, err := uc.repo.GetByID(ctx, advertID)
	if err != nil {
		return nil, advertentity.AdvertImage{},
			fmt.Errorf("get advert from repo: %w", err)
	}

	if advert.UserID != dto.UserID && !advert.IsPublic() {
		if dto.UserRole != roles.Admin {
			return nil, advertentity.AdvertImage{},
				fmt.Errorf(
					"insufficient privileges to get image of advert with status: %s: %w",
					advert.Status,
					coreerrors.ErrForbidden,
				)
		}
	}

	rc, err := uc.storage.GetByPath(advertImage.Path)
	if err != nil {
		return nil, advertentity.AdvertImage{},
			fmt.Errorf("get image by path: %w", err)
	}

	return rc, advertImage, nil
}
