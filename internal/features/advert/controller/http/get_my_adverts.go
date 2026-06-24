package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) getMyAdverts(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	userID, _, err := getUserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserID' from JWT")

		return
	}

	limit, offset, err := corehttprequest.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "get limit/offset query params")

		return
	}

	filter := advertentity.Filter{
		UserID: &userID,
	}

	dto := advertusecase.ListDTO{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
		Filter: filter,
	}

	count, adverts, err := c.useCase.List(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users adverts")

		return
	}

	response := advertsResponseFromEntities(count, adverts)

	responseHandler.JSONResponse(response, http.StatusOK)
}
