package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) getImageByID(rw http.ResponseWriter, r *http.Request) {
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

	dto := advertusecase.GetImageDTO{
		ImageID:  imageID,
		UserID:   userID,
		UserRole: userRole,
	}

	rc, image, err := c.useCase.GetImageByID(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get image by ID")

		return
	}
	defer rc.Close()

	responseHandler.FileResponse(image.Path, rc)
}
