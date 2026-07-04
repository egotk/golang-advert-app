package advertgrpc

import (
	"fmt"
	"io"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"google.golang.org/grpc"
)

func (c *Controller) GetImageByID(
	request *advertpb.GetImageByIDRequest,
	stream grpc.ServerStreamingServer[advertpb.AdvertImageResponse],
) error {
	const buffSize = 32 * 1024

	ctx := stream.Context()

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return err
	}

	dto := advertusecase.GetImageDTO{
		ImageID:  request.Id,
		UserID:   userID,
		UserRole: userRole,
	}

	rc, image, err := c.useCase.GetImageByID(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to get image by ID: %w", err)
	}
	defer rc.Close()

	response := &advertpb.AdvertImageResponse{
		ImgData: &advertpb.AdvertImageResponse_Image{
			Image: advertImageToGRPC(image),
		},
	}
	if err := stream.Send(response); err != nil {
		return fmt.Errorf("failed to send response through stream: %w", err)
	}

	buff := make([]byte, buffSize)
	for {
		n, err := rc.Read(buff)
		if n > 0 {
			response = &advertpb.AdvertImageResponse{
				ImgData: &advertpb.AdvertImageResponse_Data{
					Data: buff[:n],
				},
			}

			if err := stream.Send(response); err != nil {
				return fmt.Errorf("failed to send response through stream: %w", err)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("failed to read image from file: %w", err)
		}
	}

	return nil
}
