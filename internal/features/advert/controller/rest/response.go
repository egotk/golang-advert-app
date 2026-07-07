package advertrest

import (
	"time"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

type advertImageResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Position  int64     `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

func advertImageResponseFromEntity(i advertentity.AdvertImage) advertImageResponse {
	return advertImageResponse{
		ID:        i.ID,
		Name:      i.Name,
		Position:  i.Position,
		CreatedAt: i.CreatedAt,
	}
}

type advertImagesResponse struct {
	Count  int64                 `json:"count"`
	Images []advertImageResponse `json:"images"`
}

func imagesResponseFromEntities(images []advertentity.AdvertImage) advertImagesResponse {
	count := int64(len(images))
	imageResponses := make([]advertImageResponse, len(images))

	for i, img := range images {
		imageResponses[i] = advertImageResponseFromEntity(img)
	}

	return advertImagesResponse{
		Count:  count,
		Images: imageResponses,
	}
}

type advertResponse struct {
	ID          int64                `json:"id"`
	Version     int64                `json:"version"`
	UserID      int64                `json:"user_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Price       int64                `json:"price"`
	CategoryID  int64                `json:"category_id"`
	Status      string               `json:"status"`
	Images      advertImagesResponse `json:"advert_images"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func advertResponseFromEntity(a advertentity.Advert) advertResponse {
	return advertResponse{
		ID:          a.ID,
		Version:     a.Version,
		UserID:      a.UserID,
		Title:       a.Title,
		Description: a.Description,
		Price:       a.Price,
		CategoryID:  a.CategoryID,
		Status:      string(a.Status),
		Images:      imagesResponseFromEntities(a.Images),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

type advertsResponse struct {
	Count   int64            `json:"count"`
	Adverts []advertResponse `json:"adverts"`
}

func advertsResponseFromEntities(count int64, adverts []advertentity.Advert) advertsResponse {
	advertResponses := make([]advertResponse, len(adverts))

	for i, a := range adverts {
		advertResponses[i] = advertResponseFromEntity(a)
	}

	return advertsResponse{
		Count:   count,
		Adverts: advertResponses,
	}
}
