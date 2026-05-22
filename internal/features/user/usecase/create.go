package userusecase

import (
	"context"
	"fmt"

	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userdto "github.com/egotk/golang-advert-app/internal/features/user/usecase/dto"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCase) Create(
	ctx context.Context,
	dto userdto.Create,
) (userentity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return userentity.User{}, fmt.Errorf("failed to create bcrypt password hash: %w", err)
	}

	user := userentity.NewInitial(
		dto.Email,
		dto.FullName,
		dto.PhoneNumber,
		string(hash),
	)

	if err := uc.repo.Create(ctx, &user); err != nil {
		return userentity.User{}, fmt.Errorf("failed to store user in DB: %w", err)
	}

	return user, nil
}
