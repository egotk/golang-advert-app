package advertgrpc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"google.golang.org/grpc"
)

func (c *Controller) CreateImages(
	stream grpc.ClientStreamingServer[
		advertpb.CreateImagesRequest,
		advertpb.AdvertImagesResponse,
	],
) error {
	const maxImages = 10

	ctx := stream.Context()

	request, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to read request: %w", err)
	}

	advertID := request.GetAdvertId()
	if advertID <= 0 {
		return fmt.Errorf("'AdvertID' must be positive: %w", coreerrors.ErrInvalidArgument)
	}

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return err
	}

	var images []imageentity.Image
	count := 1
	for {
		if count > maxImages {
			return fmt.Errorf("too many images: %w", coreerrors.ErrInvalidArgument)
		}

		imageRequest, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read images: %w", err)
		}

		img := imageRequest.GetImage()
		if img == nil {
			return fmt.Errorf("failed to get image from stream: %w", coreerrors.ErrInvalidArgument)
		}

		count++

		size := len(img.Data)
		if size > imageentity.MaxImageSize {
			return fmt.Errorf("image is too big: %w", coreerrors.ErrInvalidArgument)
		}

		contentType := http.DetectContentType(img.Data)
		extension, err := imageentity.ParseExtension(contentType)
		if err != nil {
			return err
		}

		image := imageentity.Image{
			Name:      img.Name,
			Extension: extension,
			File:      bytes.NewReader(img.Data),
		}

		images = append(images, image)
	}

	dto := advertusecase.CreateImagesDTO{
		UserID:   userID,
		UserRole: userRole,
		AdvertID: advertID,
		Images:   images,
	}

	createdImages, err := c.useCase.CreateImages(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to add images to advert with id = %d: %w", advertID, err)
	}

	response := advertImagesToResponse(createdImages)

	return stream.SendAndClose(response)
}
