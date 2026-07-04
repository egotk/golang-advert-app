package advertgrpc

import (
	coregrpc "github.com/egotk/golang-advert-app/internal/core/grpc"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func filterFromRequest(f *advertpb.Filter) advertentity.Filter {
	if f == nil {
		return advertentity.Filter{}
	}

	return advertentity.Filter{
		UserID:      f.UserId,
		Title:       f.Title,
		Description: f.Description,
		MinPrice:    f.MinPrice,
		MaxPrice:    f.MaxPrice,
		CategoryID:  f.CategoryId,
		Status:      (*advertentity.Status)(f.Status),
		FromDate:    coregrpc.GRPCToTimeNullable(f.FromDate),
		ToDate:      coregrpc.GRPCToTimeNullable(f.ToDate),
		Sort:        (*advertentity.Sort)(f.Sort),
		Order:       (*advertentity.Order)(f.Order),
	}
}
