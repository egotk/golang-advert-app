package usergrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	usergrpc "github.com/egotk/golang-advert-app/internal/features/user/controller/grpc"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_Logout(t *testing.T) {
	type logoutMockBehavior func(muc *MockuseCase)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name         string
		inputRequest *userpb.LogoutRequest
		withClaims   bool
		mockBehavior logoutMockBehavior
		expectedErr  error
	}{
		{
			name:         "OK",
			inputRequest: &userpb.LogoutRequest{RefreshToken: "valid-refresh-token"},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().Logout(gomock.Any(), userusecase.LogoutDTO{
					UserID:       1,
					RefreshToken: "valid-refresh-token",
				}).Return(nil)
			},
		},
		{
			name:         "missing claims",
			inputRequest: &userpb.LogoutRequest{RefreshToken: "valid-refresh-token"},
			withClaims:   false,
			mockBehavior: func(muc *MockuseCase) {},
			expectedErr:  coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &userpb.LogoutRequest{RefreshToken: "unknown-refresh-token"},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().Logout(gomock.Any(), userusecase.LogoutDTO{
					UserID:       1,
					RefreshToken: "unknown-refresh-token",
				}).Return(coreerrors.ErrNotFound)
			},
			expectedErr: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc)

			controller := usergrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			_, err := controller.Logout(ctx, testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
		})
	}
}