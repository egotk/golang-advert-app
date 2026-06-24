package advertusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"github.com/egotk/golang-advert-app/internal/core/roles"
	"go.uber.org/zap"
)

func (uc *UseCase) Delete(
	ctx context.Context,
	dto DeleteDTO,
) error {
	log := corezaplogger.FromContext(ctx)

	advert, err := uc.repo.GetByID(ctx, dto.AdvertID)
	if err != nil {
		return fmt.Errorf("get advert: %w", err)
	}

	if dto.UserID != advert.UserID && dto.UserRole != roles.Admin {
		return fmt.Errorf(
			"user cant delete others adverts: %w",
			coreerrors.ErrForbidden,
		)
	}

	images, err := uc.repo.ListImagesByAdvertID(ctx, advert.ID)
	if err != nil {
		return fmt.Errorf("list images by advert id: %w", err)
	}

	if err := uc.repo.DeleteByID(ctx, dto.AdvertID); err != nil {
		return fmt.Errorf("delete advert: %w", err)
	}

	for _, img := range images {
		err := uc.storage.DeleteByPath(img.Path)
		if err != nil {
			log.Error("delete image", zap.String("path", img.Path), zap.Error(err))
		}
	}

	return nil
}
