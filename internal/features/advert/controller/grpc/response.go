package advertgrpc

import (
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func advertToResponse(a advertentity.Advert) *advertpb.AdvertResponse {
	return &advertpb.AdvertResponse{
		Id:          a.ID,
		Version:     a.Version,
		UserId:      a.UserID,
		Title:       a.Title,
		Description: a.Description,
		Price:       a.Price,
		CategoryId:  a.CategoryID,
		Status:      string(a.Status),
		CreatedAt:   timestamppb.New(a.CreatedAt),
		UpdatedAt:   timestamppb.New(a.UpdatedAt),
	}
}

func advertsToResponse(adverts []advertentity.Advert, count int64) *advertpb.AdvertsResponse {
	advertResponses := make([]*advertpb.AdvertResponse, len(adverts))

	for i, a := range adverts {
		advertResponses[i] = advertToResponse(a)
	}

	return &advertpb.AdvertsResponse{
		Count:   count,
		Adverts: advertResponses,
	}
}

func advertImageToGRPC(i advertentity.AdvertImage) *advertpb.AdvertImage {
	return &advertpb.AdvertImage{
		Id:        i.ID,
		Name:      i.Name,
		Position:  i.Position,
		CreatedAt: timestamppb.New(i.CreatedAt),
	}
}

func advertImagesToResponse(images []advertentity.AdvertImage) *advertpb.AdvertImagesResponse {
	count := int64(len(images))
	imageResponses := make([]*advertpb.AdvertImage, len(images))

	for i, img := range images {
		imageResponses[i] = advertImageToGRPC(img)
	}

	return &advertpb.AdvertImagesResponse{
		Count:  count,
		Images: imageResponses,
	}
}
