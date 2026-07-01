package categoryusecase

import "github.com/egotk/golang-advert-app/internal/core/nullable"

type CreateDTO struct {
	ParentID *int64
	Name     string
}

type PatchDTO struct {
	ID       int64
	ParentID nullable.Nullable[int64]
	Name     *string
}
