package userusecase

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
)

func (uc *UseCase) Logout(
	ctx context.Context,
	dto LogoutDTO,
) error {
	hash := corejwt.HashToken(dto.RefreshToken)

	if err := uc.repo.DeleteRefreshTokenByHash(ctx, dto.UserID, hash); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	return nil
}
