package advertrest

import (
	"net/http"

	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) countFavourites(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	filter, err := getListFilterQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get filter query param")

		return
	}

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := advertusecase.CountDTO{
		UserID:   userID,
		UserRole: userRole,
		Filter:   filter,
	}

	favCount, err := c.useCase.CountFavourites(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to count favourites")

		return
	}

	response := countResponse{Count: favCount}

	responseHandler.JSONResponse(response, http.StatusOK)
}
