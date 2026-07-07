package advertpostgres

import (
	"fmt"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func buildAdvertFilterConditions(filter advertentity.Filter, tableName *string, args []any) ([]string, []any) {
	conditions := []string{}

	pref := ""
	if tableName != nil {
		pref = *tableName + "."
	}

	if filter.UserID != nil {
		args = append(args, *filter.UserID)
		conditions = append(conditions, fmt.Sprintf("%suser_id = $%d", pref, len(args)))
	}

	if filter.Title != nil {
		args = append(args, "%"+*filter.Title+"%")
		conditions = append(conditions, fmt.Sprintf("%stitle ILIKE $%d", pref, len(args)))
	}

	if filter.Description != nil {
		args = append(args, "%"+*filter.Description+"%")
		conditions = append(conditions, fmt.Sprintf("%sdescription ILIKE $%d", pref, len(args)))
	}

	if filter.MinPrice != nil {
		args = append(args, *filter.MinPrice)
		conditions = append(conditions, fmt.Sprintf("%sprice >= $%d", pref, len(args)))
	}

	if filter.MaxPrice != nil {
		args = append(args, *filter.MaxPrice)
		conditions = append(conditions, fmt.Sprintf("%sprice <= $%d", pref, len(args)))
	}

	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)
		conditions = append(conditions, fmt.Sprintf("%scategory_id = $%d", pref, len(args)))
	}

	if filter.FromDate != nil {
		args = append(args, *filter.FromDate)
		conditions = append(conditions, fmt.Sprintf("%screated_at >= $%d", pref, len(args)))
	}

	if filter.ToDate != nil {
		args = append(args, *filter.ToDate)
		conditions = append(conditions, fmt.Sprintf("%screated_at <= $%d", pref, len(args)))
	}

	return conditions, args
}
