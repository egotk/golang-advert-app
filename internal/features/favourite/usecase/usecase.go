package favusecase

import "context"

type UseCase struct {
	repo repo
}

type repo interface {
	Remove(ctx context.Context, advertID int64, userID int64) error
	ListIDs(ctx context.Context, userID int64) ([]int64, error)
}

func New(repo repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
