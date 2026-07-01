package coregrpc

import (
	"context"
	"fmt"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func GetRequestIDFromMetadata(ctx context.Context) string {
	const requestIDHeader = "x-request-id"

	values := metadata.ValueFromIncomingContext(ctx, requestIDHeader)
	var requestID string
	if len(values) == 0 {
		requestID = uuid.NewString()
	} else {
		requestID = values[0]
	}

	return requestID
}

func GetTokenFromMetadata(ctx context.Context) (string, error) {
	const (
		authMDHeader = "authorization"
		bearer       = "bearer"
	)

	values := metadata.ValueFromIncomingContext(ctx, authMDHeader)
	if len(values) == 0 {
		return "", fmt.Errorf("no authorization header: %w", coreerrors.ErrUnauthorized)
	}

	authParts := strings.Split(values[0], " ")

	if len(authParts) != 2 || !strings.EqualFold(authParts[0], bearer) {
		return "", fmt.Errorf(
			"failed to parse auth header: %w",
			coreerrors.ErrUnauthorized,
		)
	}

	return authParts[1], nil
}
