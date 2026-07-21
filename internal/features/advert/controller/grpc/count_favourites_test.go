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
	"google.golang.org/protobuf/proto"
)

func TestController_CountFavourites(t *testing.T) {
	type countFavouritesMockBehavior func(muc *MockuseCase, dto advertusecase.CountDTO)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name             string
		inputRequest     *advertpb.CountRequest
		withClaims       bool
		mockBehavior     countFavouritesMockBehavior
		expectedResponse *advertpb.CountResponse
		expectedErrIs    error
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.CountRequest{},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.CountDTO) {
				muc.EXPECT().CountFavourites(gomock.Any(), dto).Return(int64(3), nil)
			},
			expectedResponse: &advertpb.CountResponse{Count: 3},
		},
		{
			name:          "missing claims",
			inputRequest:  &advertpb.CountRequest{},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.CountDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.CountRequest{},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.CountDTO) {
				muc.EXPECT().CountFavourites(gomock.Any(), dto).Return(int64(0), coreerrors.ErrInvalidArgument)
			},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.CountDTO{
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

			response, err := controller.CountFavourites(ctx, testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}