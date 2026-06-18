package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/adverts/usecase"
)

type patchRequest struct {
	Version     int     `json:"version" validate:"required,gt=0"`
	Title       *string `json:"title" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1500"`
	Price       *int    `json:"price" validate:"omitempty,gte=0"`
	CategoryID  *int    `json:"category_id" validate:"omitempty,gt=0"`
}

func (r patchRequest) toDTO(advertID, userID int, userRole string) advertusecase.PatchDTO {
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

func (c *Controller) patch(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request patchRequest
	if err := corehttprequest.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate patch advert HTTP request")

		return
	}

	advertID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'AdvertID' path param")

		return
	}

	claims, err := corejwt.ClaimsFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get claims from JWT")

		return
	}

	userID, err := claims.UserID()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserID' from claims")

		return
	}

	dto := request.toDTO(advertID, userID, claims.Role)

	advert, err := c.useCase.Patch(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch advert")

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusOK)
}
