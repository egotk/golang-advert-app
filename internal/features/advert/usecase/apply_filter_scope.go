package advertusecase

import (
	"github.com/egotk/golang-advert-app/internal/core/roles"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func applyFilterScope(
	userID int,
	userRole string,
	filter *advertentity.Filter,
) error {
	if filter.UserID != nil && userID == *filter.UserID {
		return nil
	}

	if userRole != roles.Admin {
		active := advertentity.StatusActive
		filter.Status = &active
	}

	return nil
}
