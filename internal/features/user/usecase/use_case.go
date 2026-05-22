package userusecase

import (
	"context"

	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

type UseCase struct {
	repo repo
}

type repo interface {
	Create(
		ctx context.Context,
		user *userentity.User,
	) error

	List(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]userentity.User, error)

	GetByID(
		ctx context.Context,
		id int,
	) (userentity.User, error)
}

func New(repo repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
