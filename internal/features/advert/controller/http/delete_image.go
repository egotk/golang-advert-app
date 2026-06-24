package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) deleteImage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	imageID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'ImageID' path param")

		return
	}

	userID, userRole, err := getUserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := advertusecase.DeleteImageDTO{
		UserID:   userID,
		UserRole: userRole,
		ImageID:  imageID,
	}

	if err := c.useCase.DeleteImage(ctx, dto); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete image")

		return
	}

	responseHandler.NoContentResponse()
}
