package corectx

import (
	"context"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
)

type requestIDContextKey struct{}

var requestIDKey = requestIDContextKey{}

func RequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return "", fmt.Errorf("get requestID: %w", coreerrors.ErrInvalidArgument)
	}

	return requestID, nil
}
