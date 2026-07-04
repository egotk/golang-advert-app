package advertusecase

import (
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
)

type CreateDTO struct {
	UserID      int64  `validate:"gt=0"`
	Title       string `validate:"required,min=1,max=100"`
	Description string `validate:"required,min=1,max=1500"`
	Price       int64  `validate:"gte=0"`
	CategoryID  int64  `validate:"required,gt=0"`
	Images      []imageentity.Image
}

type PatchDTO struct {
	UserID      int64
	UserRole    string
	ID          int64
	Version     int64   `validate:"required,gt=0"`
	Title       *string `validate:"omitempty,min=1,max=100"`
	Description *string `validate:"omitempty,min=1,max=1500"`
	Price       *int64  `validate:"omitempty,gte=0"`
	CategoryID  *int64  `validate:"omitempty,gt=0"`
}

type DeleteDTO struct {
	UserID   int64
	UserRole string
	AdvertID int64
}

type ArchiveDTO struct {
	UserID   int64
	UserRole string
	AdvertID int64
}

type GetByIDDTO struct {
	UserID   int64
	UserRole string
	AdvertID int64
}

type CountDTO struct {
	UserID   int64
	UserRole string
	Filter   advertentity.Filter
}

type ListDTO struct {
	UserID   int64
	UserRole string
	Limit    *int64
	Offset   *int64
	Filter   advertentity.Filter
}

type CreateImagesDTO struct {
	UserID   int64
	UserRole string
	AdvertID int64 `validate:"gt=0"`
	Images   []imageentity.Image
}

type GetImageDTO struct {
	UserID   int64
	UserRole string
	ImageID  int64
}

type DeleteImageDTO struct {
	UserID   int64
	UserRole string
	ImageID  int64
}
