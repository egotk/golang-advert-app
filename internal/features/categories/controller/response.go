package categoryhttp

import (
	categoryentity "github.com/egotk/golang-advert-app/internal/features/categories/entity"
)

type categoryResponse struct {
	ID       int    `json:"id"`
	ParentID *int   `json:"parent_id"`
	Name     string `json:"name"`
}

func categoryResponseFromEntity(c categoryentity.Category) categoryResponse {
	return categoryResponse{
		ID:       c.ID,
		ParentID: c.ParentID,
		Name:     c.Name,
	}
}

type categoriesResponse struct {
	Categories []categoryResponse `json:"categories"`
}

func categoriesResponseFromEntities(categories []categoryentity.Category) categoriesResponse {
	categoryResponses := make([]categoryResponse, len(categories))

	for i, c := range categories {
		categoryResponses[i] = categoryResponseFromEntity(c)
	}

	return categoriesResponse{Categories: categoryResponses}
}
