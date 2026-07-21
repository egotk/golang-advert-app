package advertrest

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

type patchRequest struct {
	Version     int64   `json:"version"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
	CategoryID  *int64  `json:"category_id"`
}

func (r patchRequest) toDTO(advertID, userID int64, userRole string) advertusecase.PatchDTO {
	return advertusecase.PatchDTO{
		UserID:      userID,
		UserRole:    userRole,
		ID:          advertID,
		Version:     r.Version,
		Title:       r.Title,
		Description: r.Description,
		Price:       r.Price,
		CategoryID:  r.CategoryID,
	}
}

func (c *Controller) Patch(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request patchRequest
	if err := corehttprequest.Decode(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate patch advert HTTP request")

		return
	}

	advertID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'AdvertID' path param")

		return
	}

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := request.toDTO(advertID, userID, userRole)

	advert, err := c.useCase.Patch(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch advert")

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusOK)
}
