package categoryhttp

import (
	"net/http"

	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
)

func (c *Controller) List(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	categories, err := c.useCase.List(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get categories")

		return
	}

	response := categoriesResponseFromEntities(categories)

	responseHandler.JSONResponse(response, http.StatusOK)
}
