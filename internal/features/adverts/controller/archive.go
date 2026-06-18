package adverthttp

import (
	"fmt"
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/adverts/usecase"
)

func (c *Controller) archive(rw http.ResponseWriter, r *http.Request) {
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
		responseHandler.ErrorResponse(err, "failed to get 'Claims' from JWT")

		return
	}

	userID, err := claims.UserID()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserID' from Claims")

		return
	}

	dto := advertusecase.ArchiveDTO{
		AdvertID: advertID,
		UserID:   userID,
		UserRole: claims.Role,
	}

	advert, err := c.useCase.Archive(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, fmt.Sprintf("failed to archive advert with id = %d", advertID))

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusOK)
}
