package favrest

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
)

func (c *Controller) remove(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	advertID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'AdvertID' path param")

		return
	}

	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := favusecase.RemoveDTO{
		AdvertID: advertID,
		UserID:   userID,
	}

	if err := c.useCase.Remove(ctx, dto); err != nil {
		responseHandler.ErrorResponse(err, "failed to remove advert from favourites")

		return
	}

	responseHandler.NoContentResponse()
}
