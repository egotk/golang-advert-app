package userusecase

import (
	"context"
	"fmt"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCase) CreateUser(
	ctx context.Context,
	dto CreateDTO,
) (userentity.User, error) {
	validator := corevalidator.Instance()
	if err := validator.Struct(dto); err != nil {
		return userentity.User{}, fmt.Errorf("validate DTO: %v: %w", err, coreerrors.ErrInvalidArgument)
	}

	dto.Email = strings.ToLower(dto.Email)

	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return userentity.User{}, fmt.Errorf("create bcrypt password hash: %w", err)
	}

	user := userentity.NewInitial(
		dto.Email,
		dto.FullName,
		dto.PhoneNumber,
		string(hash),
	)

	if err := uc.repo.CreateUser(ctx, &user); err != nil {
		return userentity.User{}, fmt.Errorf("store user in DB: %w", err)
	}

	return user, nil
}
