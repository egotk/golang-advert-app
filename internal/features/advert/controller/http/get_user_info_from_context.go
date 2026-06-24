package adverthttp

import (
	"context"
	"fmt"

	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
)

func getUserInfoFromContext(ctx context.Context) (int, string, error) {
	claims, err := corejwt.ClaimsFromContext(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get 'Claims' from JWT: %w", err)
	}

	userID, err := claims.UserID()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get 'UserID' from Claims: %w", err)
	}

	return userID, claims.Role, nil
}
