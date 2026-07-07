package advertentity

import (
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	"github.com/go-playground/validator/v10"
)

func StructValidations() []corevalidator.StructValidation {
	return []corevalidator.StructValidation{
		{
			Fn: func(sl validator.StructLevel) {
				filter := sl.Current().Interface().(Filter)

				if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
					sl.ReportError(filter.MinPrice, "MinPrice", "MinPrice", "ltefield", "MaxPrice")
				}

				if filter.FromDate != nil && filter.ToDate != nil && filter.FromDate.After(*filter.ToDate) {
					sl.ReportError(filter.FromDate, "FromDate", "FromDate", "ltefield", "ToDate")
				}

			},
			Types: []any{Filter{}},
		},
	}
}
