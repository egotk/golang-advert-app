package adverthttp

import (
	"time"

	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

type advertResponse struct {
	ID          int       `json:"id"`
	Version     int       `json:"version"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryID  int       `json:"category_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

type advertsResponse struct {
	Count   int              `json:"count"`
	Adverts []advertResponse `json:"adverts"`
}

func advertsResponseFromEntities(count int, adverts []advertentity.Advert) advertsResponse {
	advertResponses := make([]advertResponse, len(adverts))

	for i, a := range adverts {
		advertResponses[i] = advertResponseFromEntity(a)
	}

	return advertsResponse{
		Count:   count,
		Adverts: advertResponses,
	}
}
