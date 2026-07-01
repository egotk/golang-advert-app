package adverthttp

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
)

type createImagesRequest struct {
	AdvertID      int64 `validate:"required,gt=0"`
	ImagesHeaders []*multipart.FileHeader
}

// TODO: вынести общую с createRequest логику
func createImagesRequestFromMultipart(r *http.Request) (createImagesRequest, error) {
	const (
		advertIDKey  = "advert_id"
		imagesKey    = "images"
		maxImageSize = 10 << 20
		maxMemory    = 8 << 20
	)

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return createImagesRequest{}, fmt.Errorf(
			"bad multipart request: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	advertID, err := strconv.ParseInt(r.FormValue(advertIDKey), 10, 64)
	if err != nil {
		return createImagesRequest{}, fmt.Errorf(
			"get 'AdvertID' from multipart: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	request := createImagesRequest{AdvertID: advertID}

	if err := corevalidator.Instance().Struct(request); err != nil {
		return createImagesRequest{}, fmt.Errorf(
			"validate request: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	headers := r.MultipartForm.File[imagesKey]

	for _, h := range headers {
		if h.Size > maxImageSize {
			return createImagesRequest{}, fmt.Errorf(
				"image is too big: %w",
				coreerrors.ErrInvalidArgument,
			)
		}
	}

	request.ImagesHeaders = headers

	return request, nil
}

func (c *Controller) createImages(rw http.ResponseWriter, r *http.Request) {
	const contentTypeSeekLen = 512

	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	request, err := createImagesRequestFromMultipart(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate create images HTTP request")

		return
	}

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user info from JWT")

		return
	}

	var images []imageentity.Image

	for _, h := range request.ImagesHeaders {

		file, err := h.Open()
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to parse images")

			return
		}
		defer file.Close()

		buff := make([]byte, contentTypeSeekLen)
		n, err := io.ReadFull(file, buff)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			responseHandler.ErrorResponse(
				err,
				fmt.Sprintf("failed to detect image %s content type", h.Filename),
			)

			return
		}

		contentType := http.DetectContentType(buff[:n])

		extension, err := imageentity.ParseExtension(contentType)
		if err != nil {
			responseHandler.ErrorResponse(
				err,
				fmt.Sprintf("failed to validate image %s", h.Filename),
			)

			return
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			responseHandler.ErrorResponse(
				err,
				"failed to seek start of image",
			)

			return
		}

		image := imageentity.Image{
			Name:      h.Filename,
			Extension: extension,
			File:      file,
		}

		images = append(images, image)
	}

	dto := advertusecase.CreateImagesDTO{
		UserID:   userID,
		UserRole: userRole,
		AdvertID: request.AdvertID,
		Images:   images,
	}

	createdImages, err := c.useCase.CreateImages(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			fmt.Sprintf("failed to add images to advert with id=%d", dto.AdvertID),
		)

		return
	}

	response := imagesResponseFromEntities(createdImages)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
