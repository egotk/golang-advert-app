package categoryhttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
)

func (c *Controller) Delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	categoryID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'ID' path param")

		return
	}

	if err := c.useCase.Delete(ctx, categoryID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete category")

		return
	}

	responseHandler.NoContentResponse()
}
