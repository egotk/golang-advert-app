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

func TestController_Delete(t *testing.T) {
	type deleteMockBehavior func(muc *MockuseCase, dto advertusecase.DeleteDTO)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name          string
		inputRequest  *advertpb.DeleteRequest
		withClaims    bool
		mockBehavior  deleteMockBehavior
		expectedErrIs error
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.DeleteRequest{Id: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.DeleteDTO) {
				muc.EXPECT().Delete(gomock.Any(), dto).Return(nil)
			},
		},
		{
			name:          "missing claims",
			inputRequest:  &advertpb.DeleteRequest{Id: 1},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.DeleteDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.DeleteRequest{Id: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.DeleteDTO) {
				muc.EXPECT().Delete(gomock.Any(), dto).Return(coreerrors.ErrForbidden)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.DeleteDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: testCase.inputRequest.Id,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			_, err := controller.Delete(ctx, testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}