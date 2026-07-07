package advertpostgres

import (
	"fmt"

	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func orderBy(
	sort *advertentity.Sort,
	tableName *string,
	order *advertentity.Order,
) string {
	pref := ""
	if tableName != nil {
		pref = *tableName + "."
	}

	if sort == nil {
		return fmt.Sprintf(" ORDER BY %sid ASC", pref)
	}

	var column string
	switch *sort {
	case advertentity.SortByPrice:
		column = fmt.Sprintf("%sprice", pref)
	case advertentity.SortByViews:
		column = fmt.Sprintf("%sviews_count", pref)
	case advertentity.SortByDate:
		column = fmt.Sprintf("%screated_at", pref)
	default:
		return fmt.Sprintf(" ORDER BY %sid ASC", pref)
	}

	direction := "ASC"
	if order != nil && *order == advertentity.OrderDesc {
		direction = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s, %sid ASC", column, direction, pref)
}
