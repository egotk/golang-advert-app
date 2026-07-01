package userusecase

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
)

func (uc *UseCase) Logout(
	ctx context.Context,
	dto LogoutDTO,
) error {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	hash := corejwt.HashToken(dto.RefreshToken)

	if err := uc.repo.DeleteRefreshTokenByHash(ctx, dto.UserID, hash); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	return nil
}
