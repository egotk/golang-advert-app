package userhttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
)

type refreshTokensRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1,max=512"`
}

func (r refreshTokensRequest) toDTO() userusecase.RefreshTokensDTO {
	return userusecase.RefreshTokensDTO{
		RefreshToken: r.RefreshToken,
	}
}

type refreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Controller) refreshTokens(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request refreshTokensRequest
	if err := corehttprequest.Decode(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")

		return
	}

	dto := request.toDTO()

	tokens, err := c.useCase.RefreshTokens(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh tokens")

		return
	}

	response := refreshTokensResponse{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
