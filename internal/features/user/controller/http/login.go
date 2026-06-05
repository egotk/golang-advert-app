package userhttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
)

type loginRequest struct {
	Email    string `json:"email"    validate:"required,min=3,max=255,email"`
	Password string `json:"password" validate:"required,min=8,max=36"`
}

func (r loginRequest) toDTO() userusecase.LoginDTO {
	return userusecase.LoginDTO{
		Email:    r.Email,
		Password: r.Password,
	}
}

type loginResponse struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Controller) login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request loginRequest
	if err := corehttprequest.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	dto := request.toDTO()

	result, err := c.useCase.Login(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")

		return
	}

	response := loginResponse{
		UserID:       result.UserID,
		AccessToken:  result.Tokens.Access,
		RefreshToken: result.Tokens.Refresh,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
