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
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
)

type createRequest struct {
	Title         string
	Description   string
	Price         int64
	CategoryID    int64
	ImagesHeaders []*multipart.FileHeader
}

func (r createRequest) toDTO(
	userID int64,
	images []imageentity.Image,
) advertusecase.CreateDTO {
	return advertusecase.CreateDTO{
		UserID:      userID,
		Title:       r.Title,
		Description: r.Description,
		Price:       r.Price,
		CategoryID:  r.CategoryID,
		Images:      images,
	}
}

// TODO: вынести общую с addImageRequest логику
func createRequestFromMultipart(r *http.Request) (createRequest, error) {
	const (
		titleKey        = "title"
		descriptionKey  = "description"
		priceKey        = "price"
		categoryIDKey   = "category_id"
		imagesKey       = "images"
		maxAdvertImages = 10
		maxMemory       = 8 * 1024 * 1024
	)

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return createRequest{}, fmt.Errorf(
			"bad multipart request: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	title := r.FormValue(titleKey)

	description := r.FormValue(descriptionKey)

	price, err := strconv.ParseInt(r.FormValue(priceKey), 10, 64)
	if err != nil {
		return createRequest{}, fmt.Errorf(
			"get 'Price' from multipart: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	categoryID, err := strconv.ParseInt(r.FormValue(categoryIDKey), 10, 64)
	if err != nil {
		return createRequest{}, fmt.Errorf(
			"get 'CategoryID' from multipart: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	request := createRequest{
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
	}

	headers := r.MultipartForm.File[imagesKey]
	if len(headers) > maxAdvertImages {
		return createRequest{}, fmt.Errorf(
			"too many images: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	for _, h := range headers {
		if h.Size > imageentity.MaxImageSize {
			return createRequest{}, fmt.Errorf(
				"image is too big: %w",
				coreerrors.ErrInvalidArgument,
			)
		}
	}

	request.ImagesHeaders = headers

	return request, nil
}

func (c *Controller) create(rw http.ResponseWriter, r *http.Request) {
	const contentTypeSeekLen = 512

	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	request, err := createRequestFromMultipart(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate create advert HTTP request")

		return
	}

	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

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

	dto := request.toDTO(userID, images)

	advert, err := c.useCase.Create(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create advert")

		return
	}

	response := advertResponseFromEntity(advert)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
