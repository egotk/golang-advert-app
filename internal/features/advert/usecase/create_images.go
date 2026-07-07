package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	"go.uber.org/zap"
)

func (uc *UseCase) CreateImages(ctx context.Context, dto CreateImagesDTO) (_ []advertentity.AdvertImage, err error) {
	const maxImagesLen = 10
	var savedPaths []string
	defer func() {
		if err != nil {
			log := corezaplogger.FromContext(ctx)

			for _, p := range savedPaths {
				if err := uc.storage.DeleteByPath(p); err != nil {
					log.Error(
						"delete saved images after error",
						zap.String("path", p),
						zap.Error(err),
					)
				}
			}
		}
	}()

	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return nil, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	if len(dto.Images) == 0 {
		return nil, fmt.Errorf(
			"empty add images request: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	advert, err := uc.repo.GetByID(ctx, dto.AdvertID)
	if err != nil {
		return nil, fmt.Errorf("get adverts from repo: %w", err)
	}

	if len(advert.Images)+len(dto.Images) > 10 {
		return nil, fmt.Errorf(
			"advert images cant exceed %d items: %w",
			maxImagesLen,
			coreerrors.ErrInvalidArgument,
		)
	}

	if advert.UserID != dto.UserID {
		if dto.UserRole != roles.Admin {
			return nil, fmt.Errorf(
				"insufficient privileges to delete image: %w",
				coreerrors.ErrForbidden,
			)
		}
	}

	var images []advertentity.AdvertImage
	for _, i := range dto.Images {
		path, err := uc.storage.Save(i.Extension, i.File)
		if err != nil {
			return nil, fmt.Errorf("save image: %w", err)
		}

		savedPaths = append(savedPaths, path)

		image := advertentity.NewAdvertImageUninitialized(i.Name, path)

		images = append(images, image)
	}

	if err := uc.repo.CreateImages(ctx, dto.AdvertID, images); err != nil {
		return nil, fmt.Errorf("store images in repo: %w", err)
	}

	return images, nil
}
