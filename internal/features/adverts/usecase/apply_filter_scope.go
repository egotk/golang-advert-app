package advertusecase

import (
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/adverts/entity"
)

func applyFilterScope(
	userRole string,
	filter *advertentity.Filter,
) error {
	if userRole != roles.Admin {
		active := advertentity.StatusActive
		filter.Status = &active
	}

	return nil
}
