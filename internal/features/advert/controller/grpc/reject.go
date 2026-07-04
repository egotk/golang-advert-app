package advertgrpc

import (
	"context"
	"fmt"

	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) Reject(
	ctx context.Context,
	request *advertpb.RejectRequest,
) (*advertpb.AdvertResponse, error) {
	advert, err := c.useCase.Reject(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to reject advert with id=%d: %w", request.Id, err)
	}

	response := advertToResponse(advert)

	return response, nil
}
