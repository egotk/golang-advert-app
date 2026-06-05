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

func (c Claims) UserID() (int, error) {
	userID, err := strconv.Atoi(c.Subject)
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
