package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/adverts/usecase"
)

func (c *Controller) delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	advertID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'AdvertID' path param")

		return
	}

	claims, err := corejwt.ClaimsFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get claims from JWT")

		return
	}

	userID, err := claims.UserID()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserID' from claims")

		return
	}

	dto := advertusecase.DeleteDTO{
		UserID:   userID,
		UserRole: claims.Role,
		AdvertID: advertID,
	}

	if err := c.useCase.Delete(ctx, dto); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete advert")

		return
	}

	responseHandler.NoContentResponse()
}
