package advertgrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertgrpc "github.com/egotk/golang-advert-app/internal/features/advert/controller/grpc"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_AddToFavourites(t *testing.T) {
	type addToFavouritesMockBehavior func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name          string
		inputRequest  *advertpb.AddToFavouritesRequest
		withClaims    bool
		mockBehavior  addToFavouritesMockBehavior
		expectedErrIs error
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.AddToFavouritesRequest{AdvertId: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {
				muc.EXPECT().AddToFavourites(gomock.Any(), dto).Return(nil)
			},
		},
		{
			name:          "missing claims",
			inputRequest:  &advertpb.AddToFavouritesRequest{AdvertId: 1},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.AddToFavouritesRequest{AdvertId: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {
				muc.EXPECT().AddToFavourites(gomock.Any(), dto).Return(coreerrors.ErrForbidden)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.AddToFavouritesDTO{
				AdvertID: testCase.inputRequest.AdvertId,
				UserID:   1,
				UserRole: "user",
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			_, err := controller.AddToFavourites(ctx, testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}