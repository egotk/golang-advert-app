package advertpostgres

import (
	"fmt"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func buildFilterConditions(filter advertentity.Filter) ([]string, []any) {
	conditions := []string{}
	args := []any{}

	if filter.UserID != nil {
		args = append(args, *filter.UserID)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if filter.Title != nil {
		args = append(args, "%"+*filter.Title+"%")
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", len(args)))
	}

	if filter.Description != nil {
		args = append(args, "%"+*filter.Description+"%")
		conditions = append(conditions, fmt.Sprintf("description ILIKE $%d", len(args)))
	}

	if filter.MinPrice != nil {
		args = append(args, *filter.MinPrice)
		conditions = append(conditions, fmt.Sprintf("price >= $%d", len(args)))
	}

	if filter.MaxPrice != nil {
		args = append(args, *filter.MaxPrice)
		conditions = append(conditions, fmt.Sprintf("price <= $%d", len(args)))
	}

	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
	}

	if filter.FromDate != nil {
		args = append(args, *filter.FromDate)
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)))
	}

	if filter.ToDate != nil {
		args = append(args, *filter.ToDate)
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", len(args)))
	}

	return conditions, args
}
