package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	"go.uber.org/zap"
)

func (uc *UseCase) Create(
	ctx context.Context,
	dto CreateDTO,
) (_ advertentity.Advert, err error) {
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
		return advertentity.Advert{}, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	advert := advertentity.NewInitial(
		dto.UserID,
		dto.Title,
		dto.Description,
		dto.Price,
		dto.CategoryID,
	)

	var images []advertentity.AdvertImage
	for idx, i := range dto.Images {
		path, err := uc.storage.Save(i.Extension, i.File)
		if err != nil {
			return advertentity.Advert{}, fmt.Errorf("save image: %w", err)
		}

		savedPaths = append(savedPaths, path)

		image := advertentity.NewAdvertImageUninitialized(
			dto.Images[idx].Name,
			path,
		)

		images = append(images, image)
	}

	advert.Images = images

	if err := uc.repo.Create(ctx, &advert); err != nil {
		return advertentity.Advert{}, fmt.Errorf("store advert in DB: %w", err)
	}

	return advert, nil
}
