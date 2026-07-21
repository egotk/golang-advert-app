package advertrest

import (
	"fmt"
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
)

func (c *Controller) Reject(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	advertID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'AdvertID' path param")

		return
	}

	advert, err := c.useCase.Reject(ctx, advertID)
	if err != nil {
		responseHandler.ErrorResponse(err, fmt.Sprintf("failed to reject advert with id = %d", advertID))

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusOK)
}
