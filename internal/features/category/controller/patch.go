package categoryhttp

import (
	"fmt"
	"net/http"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"github.com/egotk/golang-advert-app/internal/core/nullable"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
)

type patchRequest struct {
	ParentID nullable.Nullable[int64] `json:"parent_id"`
	Name     *string                  `json:"name"`
}

func (r patchRequest) toDTO(id int64) categoryusecase.PatchDTO {
	return categoryusecase.PatchDTO{
		ID:       id,
		ParentID: r.ParentID,
		Name:     r.Name,
	}
}

func (r patchRequest) Validate() error {
	if r.ParentID.Set && r.ParentID.Value != nil {
		if *r.ParentID.Value < 1 {
			return fmt.Errorf("'ParentID' must be positive: %w", coreerrors.ErrInvalidArgument)
		}
	}

	if r.Name != nil {
		nameLen := len([]rune(*r.Name))

		if nameLen < 1 || nameLen > 100 {
			return fmt.Errorf("'Name' must be between 1 and 100 characters: %w", coreerrors.ErrInvalidArgument)
		}
	}

	return nil
}

func (c *Controller) patch(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request patchRequest
	if err := corehttprequest.Decode(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate patch category HTTP request")

		return
	}

	categoryID, err := corehttprequest.GetIntPathParam("id", r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'ID' path param")

		return
	}

	dto := request.toDTO(categoryID)

	category, err := c.useCase.Patch(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch category")

		return
	}

	response := categoryResponseFromEntity(category)

	responseHandler.JSONResponse(response, http.StatusOK)
}
