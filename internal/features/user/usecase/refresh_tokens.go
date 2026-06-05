package userusecase

import (
	"context"
	"fmt"
	"time"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func (uc *UseCase) RefreshTokens(
	ctx context.Context,
	dto RefreshTokensDTO,
) (TokensDTO, error) {
	oldHash := corejwt.HashToken(dto.RefreshToken)

	oldToken, err := uc.repo.GetRefreshTokenByHash(ctx, oldHash)
	if err != nil {
		return TokensDTO{}, fmt.Errorf("get refresh token: %w", err)
	}

	if oldToken.ExpiresAt.Before(time.Now()) {
		return TokensDTO{}, userentity.ErrRefreshTokenExpired
	}

	user, err := uc.repo.GetUserByID(ctx, oldToken.UserID)
	if err != nil {
		return TokensDTO{}, fmt.Errorf("get user: %w", err)
	}

	newPair, err := uc.jwtService.IssuePair(user.Role, oldToken.UserID)
	if err != nil {
		return TokensDTO{}, fmt.Errorf("issue jwt pair: %w", err)
	}

	newHash := corejwt.HashToken(newPair.RefreshToken.Token)

	newRefresh := userentity.RefreshToken{
		Hash:      newHash,
		UserID:    oldToken.UserID,
		IssuedAt:  newPair.RefreshToken.IssuedAt,
		ExpiresAt: newPair.RefreshToken.ExpiresAt,
	}

	if err := uc.repo.ReissueRefreshToken(ctx, oldToken.UserID, oldHash, newRefresh); err != nil {
		return TokensDTO{}, fmt.Errorf("reissue refresh token: %w", err)
	}

	tokensDTO := TokensDTO{
		Access:  newPair.AccessToken,
		Refresh: newPair.RefreshToken.Token,
	}

	return tokensDTO, nil
}
