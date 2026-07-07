package favrest

import (
	"net/http"

	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
)

type listIDsResponse struct {
	FavouriteIDs []int64 `json:"ids"`
}

func (c *Controller) listIDs(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	ids, err := c.useCase.ListIDs(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list favourite adverts ids")

		return
	}

	response := listIDsResponse{FavouriteIDs: ids}

	responseHandler.JSONResponse(response, http.StatusOK)
}
