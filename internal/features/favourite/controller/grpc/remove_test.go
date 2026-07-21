package favgrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	favgrpc "github.com/egotk/golang-advert-app/internal/features/favourite/controller/grpc"
	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
	favpb "github.com/egotk/golang-advert-app/internal/gen/favourite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_Remove(t *testing.T) {
	type removeMockBehavior func(muc *MockuseCase, dto favusecase.RemoveDTO)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name         string
		inputRequest *favpb.RemoveRequest
		withClaims   bool
		mockBehavior removeMockBehavior
		expectedErr  error
	}{
		{
			name:         "OK",
			inputRequest: &favpb.RemoveRequest{AdvertId: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto favusecase.RemoveDTO) {
				muc.EXPECT().Remove(gomock.Any(), dto).Return(nil)
			},
		},
		{
			name:         "missing claims",
			inputRequest: &favpb.RemoveRequest{AdvertId: 1},
			withClaims:   false,
			mockBehavior: func(muc *MockuseCase, dto favusecase.RemoveDTO) {},
			expectedErr:  coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &favpb.RemoveRequest{AdvertId: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto favusecase.RemoveDTO) {
				muc.EXPECT().Remove(gomock.Any(), dto).Return(coreerrors.ErrNotFound)
			},
			expectedErr: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			dto := favusecase.RemoveDTO{AdvertID: testCase.inputRequest.AdvertId, UserID: 1}
			testCase.mockBehavior(muc, dto)

			controller := favgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			_, err := controller.Remove(ctx, testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
		})
	}
}
