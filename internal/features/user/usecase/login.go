package userusecase

import (
	"context"
	"fmt"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCase) Login(
	ctx context.Context,
	dto LoginDTO,
) (LoginResultDTO, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return LoginResultDTO{}, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	user, err := uc.repo.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return LoginResultDTO{}, fmt.Errorf("get user with email: %w", err)
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return LoginResultDTO{}, userentity.ErrUserIsLocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		lockedUntil, err := uc.repo.IncrementFailedLoginCount(ctx, user.ID, user.Version)
		if err != nil {
			return LoginResultDTO{}, fmt.Errorf("increment 'FailedLoginCount': %w", err)
		}
		if lockedUntil != nil {
			return LoginResultDTO{}, fmt.Errorf("%w until %s", userentity.ErrUserIsLocked, lockedUntil)
		}

		return LoginResultDTO{}, userentity.ErrInvalidPassword
	}

	if err := uc.repo.ResetFailedLoginCount(ctx, user.ID, user.Version); err != nil {
		return LoginResultDTO{}, fmt.Errorf("reset 'FailedLoginCount': %w", err)
	}

	jwtPair, err := uc.jwtService.IssuePair(user.Role, user.ID)
	if err != nil {
		return LoginResultDTO{}, fmt.Errorf("issue jwt pair: %w", err)
	}

	refreshHash := corejwt.HashToken(jwtPair.RefreshToken.Token)

	refresh := userentity.RefreshToken{
		UserID:    user.ID,
		Hash:      refreshHash,
		IssuedAt:  jwtPair.RefreshToken.IssuedAt,
		ExpiresAt: jwtPair.RefreshToken.ExpiresAt,
	}

	if err := uc.repo.CreateRefreshToken(ctx, refresh); err != nil {
		return LoginResultDTO{}, fmt.Errorf("store refresh token: %w", err)
	}

	result := LoginResultDTO{
		UserID: user.ID,
		Tokens: TokensDTO{
			Access:  jwtPair.AccessToken,
			Refresh: jwtPair.RefreshToken.Token,
		},
	}

	return result, nil
}
