package adverthttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/adverts/usecase"
)

type createRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1,max=1500"`
	Price       int    `json:"price" validate:"gte=0"`
	CategoryID  int    `json:"category_id" validate:"required,gt=0"`
}

func (r createRequest) toDTO(userID int) advertusecase.CreateDTO {
	return advertusecase.CreateDTO{
		UserID:      userID,
		Title:       r.Title,
		Description: r.Description,
		Price:       r.Price,
		CategoryID:  r.CategoryID,
	}
}

func (c *Controller) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request createRequest
	if err := corehttprequest.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate create advert HTTP request")

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

	dto := request.toDTO(userID)

	advert, err := c.useCase.Create(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create advert")

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
