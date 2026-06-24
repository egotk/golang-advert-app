package categoryusecase

import "github.com/egotk/golang-advert-app/internal/core/nullable"

type CreateDTO struct {
	ParentID *int
	Name     string
}

type PatchDTO struct {
	ID       int
	ParentID nullable.Nullable[int]
	Name     *string
}
