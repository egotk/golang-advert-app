package userusecase

import (
	"context"
	"fmt"

	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCase) CreateUser(
	ctx context.Context,
	dto CreateDTO,
) (userentity.User, error) {
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
