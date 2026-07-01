package corejwt

import (
	"context"
	"fmt"
	"strconv"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

type claimsContextKey struct{}

var key = claimsContextKey{}

type Claims struct {
	Role string `json:"role"`

	jwt.RegisteredClaims
}

func (c Claims) UserID() (int64, error) {
	userID, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid 'sub': %s: %w", c.Subject, coreerrors.ErrInvalidArgument)
	}

	return userID, nil
}

func ClaimsToContext(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, key, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, error) {
	claims, ok := ctx.Value(key).(Claims)
	if !ok {
		return Claims{}, fmt.Errorf("get token claims: %w", coreerrors.ErrUnauthorized)
	}

	return claims, nil
}

func UserInfoFromContext(ctx context.Context) (int64, string, error) {
	claims, err := ClaimsFromContext(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get 'Claims' from JWT: %w", err)
	}

	userID, err := claims.UserID()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get 'UserID' from Claims: %w", err)
	}

	return userID, claims.Role, nil
}
