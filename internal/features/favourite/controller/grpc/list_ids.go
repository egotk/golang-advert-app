package favgrpc

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	favpb "github.com/egotk/golang-advert-app/internal/gen/favourite"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) ListIDs(ctx context.Context, _ *emptypb.Empty) (*favpb.ListIDsResponse, error) {
	userID, _, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get 'UserInfo' from JWT: %w", err)
	}

	ids, err := c.useCase.ListIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list favourite adverts ids: %w", err)
	}

	response := &favpb.ListIDsResponse{
		Ids: ids,
	}

	return response, nil
}
