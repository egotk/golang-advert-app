package userhttp

import (
	"net/http"
	"strings"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
)

type createRequest struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (r createRequest) toDTO() userusecase.CreateDTO {
	return userusecase.CreateDTO{
		Email:       strings.ToLower(r.Email),
		FullName:    r.FullName,
		PhoneNumber: r.PhoneNumber,
		Password:    r.Password,
	}
}

func (r createRequest) Validate() error {
	if err := userentity.ValidateEmail(r.Email); err != nil {
		return err
	}

	if err := userentity.ValidateFullName(r.FullName); err != nil {
		return err
	}

	if err := userentity.ValidatePhoneNumber(r.PhoneNumber); err != nil {
		return err
	}

	return userentity.ValidatePassword(r.Password)
}

type createResponse dtoResponse

func (c *Controller) createUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	var request createRequest
	if err := corehttprequest.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	dto := request.toDTO()

	user, err := c.useCase.CreateUser(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := createResponse(dtoResponseFromEntity(user))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
