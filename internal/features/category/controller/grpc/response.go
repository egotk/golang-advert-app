package categorygrpc

import (
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
)

func categoryToResponse(c categoryentity.Category) *categorypb.CategoryResponse {
	return &categorypb.CategoryResponse{
		Id:       c.ID,
		ParentId: c.ParentID,
		Name:     c.Name,
	}
}

func categoriesToResponse(categories []categoryentity.Category) *categorypb.CategoriesResponse {
	categoryResponses := make([]*categorypb.CategoryResponse, len(categories))

	for i, c := range categories {
		categoryResponses[i] = categoryToResponse(c)
	}

	return &categorypb.CategoriesResponse{
		Categories: categoryResponses,
	}
}
