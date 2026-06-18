package categoryhttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/categories/usecase"
)

type createRequest struct {
	ParentID *int   `json:"parent_id" validate:"omitempty,gte=1"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
}

func (r createRequest) toDTO() categoryusecase.CreateDTO {
	return categoryusecase.CreateDTO{
		ParentID: r.ParentID,
		Name:     r.Name,
	}
}

func (c *Controller) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request createRequest
	if err := corehttprequest.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate create category HTTP request")

		return
	}

	dto := request.toDTO()

	category, err := c.useCase.Create(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create category")

		return
	}

	response := categoryResponseFromEntity(category)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
