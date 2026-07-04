package categoryusecase

import "github.com/egotk/golang-advert-app/internal/core/nullable"

type CreateDTO struct {
	ParentID *int64 `validate:"omitempty,gt=0"`
	Name     string `validate:"required,min=1,max=100"`
}

type PatchDTO struct {
	ID       int64 `validate:"gt=0"`
	ParentID nullable.Nullable[int64]
	Name     *string `validate:"omitempty,min=1,max=100"`
}
