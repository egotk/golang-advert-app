package advertrest

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) listFavourites(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	limit, offset, err := corehttprequest.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset query param")

		return
	}

	filter, err := getListFilterQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get list filter query params")

		return
	}

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := advertusecase.ListDTO{
		UserID:   userID,
		UserRole: userRole,
		Limit:    limit,
		Offset:   offset,
		Filter:   filter,
	}

	count, adverts, err := c.useCase.ListFavourites(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get favourite adverts")

		return
	}

	response := advertsResponseFromEntities(count, adverts)

	responseHandler.JSONResponse(response, http.StatusOK)
}
