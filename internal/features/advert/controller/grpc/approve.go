package advertgrpc

import (
	"context"
	"fmt"

	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
)

func (c *Controller) Approve(
	ctx context.Context,
	request *advertpb.ApproveRequest,
) (*advertpb.AdvertResponse, error) {
	advert, err := c.useCase.Approve(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to approve advert with id = %d: %w", request.Id, err)
	}

	response := advertToResponse(advert)

	return response, nil
}
