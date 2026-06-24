package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) getByID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	advertID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'AdvertID' path param")

		return
	}

	userID, userRole, err := getUserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := advertusecase.GetByIDDTO{
		UserID:   userID,
		UserRole: userRole,
		AdvertID: advertID,
	}

	advert, err := c.useCase.GetByID(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get advert by ID")

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusOK)
}
