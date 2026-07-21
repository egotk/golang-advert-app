package userhttp

import (
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
)

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r logoutRequest) toDTO(userID int64) userusecase.LogoutDTO {
	return userusecase.LogoutDTO{
		UserID:       userID,
		RefreshToken: r.RefreshToken,
	}
}

func (c *Controller) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request logoutRequest
	if err := corehttprequest.Decode(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")

		return
	}

	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := request.toDTO(userID)

	if err := c.useCase.Logout(ctx, dto); err != nil {
		responseHandler.ErrorResponse(err, "failed to logout")

		return
	}

	responseHandler.NoContentResponse()
}
