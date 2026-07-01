package advertusecase

import (
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
)

type CreateDTO struct {
	UserID      int64
	Title       string
	Description string
	Price       int64
	CategoryID  int64
	Images      []imageentity.Image
}

type PatchDTO struct {
	UserID      int64
	UserRole    string
	ID          int64
	Version     int64
	Title       *string
	Description *string
	Price       *int64
	CategoryID  *int64
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
	AdvertID int64
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
