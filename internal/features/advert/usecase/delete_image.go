package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	"go.uber.org/zap"
)

func (uc *UseCase) DeleteImage(
	ctx context.Context,
	dto DeleteImageDTO,
) error {
	if dto.ImageID <= 0 {
		return fmt.Errorf("'ImageID' must be positive")
	}

	log := corezaplogger.FromContext(ctx)

	advertID, image, err := uc.repo.GetImageByID(ctx, dto.ImageID)
	if err != nil {
		return fmt.Errorf("get image from repo: %w", err)
	}

	advert, err := uc.repo.GetByID(ctx, advertID)
	if err != nil {
		return fmt.Errorf("get advert from repo: %w", err)
	}

	if advert.UserID != dto.UserID {
		if dto.UserRole != roles.Admin {
			return fmt.Errorf(
				"insufficient privileges to delete image: %w",
				coreerrors.ErrForbidden,
			)
		}
	}

	if err := uc.repo.DeleteImageByID(ctx, dto.ImageID); err != nil {
		return fmt.Errorf("delete image from repo: %w", err)
	}

	if err := uc.storage.DeleteByPath(image.Path); err != nil {
		log.Error("delete image from local storage", zap.String("path", image.Path), zap.Error(err))
	}

	return nil
}
