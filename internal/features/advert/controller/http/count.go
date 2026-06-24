package adverthttp

import (
	"net/http"

	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

type countResponse struct {
	Count int `json:"count"`
}

func (c *Controller) count(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

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

	dto := advertusecase.CountDTO{
		UserID:   userID,
		UserRole: userRole,
		Filter:   filter,
	}

	count, err := c.useCase.Count(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to count adverts")

		return
	}

	response := countResponse{Count: count}

	responseHandler.JSONResponse(response, http.StatusOK)
}
