package advertusecase

import (
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
)

type CreateDTO struct {
	UserID      int
	Title       string
	Description string
	Price       int
	CategoryID  int
	Images      []imageentity.Image
}

type PatchDTO struct {
	UserID      int
	UserRole    string
	ID          int
	Version     int
	Title       *string
	Description *string
	Price       *int
	CategoryID  *int
}

type DeleteDTO struct {
	UserID   int
	UserRole string
	AdvertID int
}

type ArchiveDTO struct {
	UserID   int
	UserRole string
	AdvertID int
}

type GetByIDDTO struct {
	UserID   int
	UserRole string
	AdvertID int
}

type CountDTO struct {
	UserID   int
	UserRole string
	Filter   advertentity.Filter
}

type ListDTO struct {
	UserID   int
	UserRole string
	Limit    *int
	Offset   *int
	Filter   advertentity.Filter
}

type CreateImagesDTO struct {
	UserID   int
	UserRole string
	AdvertID int
	Images   []imageentity.Image
}

type GetImageDTO struct {
	UserID   int
	UserRole string
	ImageID  int
}

type DeleteImageDTO struct {
	UserID   int
	UserRole string
	ImageID  int
}
