package advertusecase

import advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"

type CreateDTO struct {
	UserID      int
	Title       string
	Description string
	Price       int
	CategoryID  int
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
	AdvertID int
	UserID   int
	UserRole string
}

type GetByIDDTO struct {
	AdvertID int
	UserID   int
	UserRole string
}

type CountDTO struct {
	Filter   advertentity.Filter
	UserRole string
}

type ListDTO struct {
	Limit    *int
	Offset   *int
	Filter   advertentity.Filter
	UserRole string
}
