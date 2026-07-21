package usergrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	usergrpc "github.com/egotk/golang-advert-app/internal/features/user/controller/grpc"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestController_RefreshTokens(t *testing.T) {
	type refreshTokensMockBehavior func(muc *MockuseCase, dto userusecase.RefreshTokensDTO)

	testTable := []struct {
		name             string
		inputRequest     *userpb.RefreshTokensRequest
		inputDTO         userusecase.RefreshTokensDTO
		mockBehavior     refreshTokensMockBehavior
		expectedResponse *userpb.RefreshTokensResponse
		expectedErr      error
	}{
		{
			name:         "OK",
			inputRequest: &userpb.RefreshTokensRequest{RefreshToken: "old-refresh-token"},
			inputDTO:     userusecase.RefreshTokensDTO{RefreshToken: "old-refresh-token"},
			mockBehavior: func(muc *MockuseCase, dto userusecase.RefreshTokensDTO) {
				muc.EXPECT().RefreshTokens(gomock.Any(), dto).Return(
					userusecase.TokensDTO{
						Access:  "new-access-token",
						Refresh: "new-refresh-token",
					}, nil)
			},
			expectedResponse: &userpb.RefreshTokensResponse{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			},
		},
		{
			name:         "usecase error",
			inputRequest: &userpb.RefreshTokensRequest{RefreshToken: "expired-refresh-token"},
			inputDTO:     userusecase.RefreshTokensDTO{RefreshToken: "expired-refresh-token"},
			mockBehavior: func(muc *MockuseCase, dto userusecase.RefreshTokensDTO) {
				muc.EXPECT().RefreshTokens(gomock.Any(), dto).Return(userusecase.TokensDTO{}, coreerrors.ErrUnauthorized)
			},
			expectedErr: coreerrors.ErrUnauthorized,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := usergrpc.New(muc)

			response, err := controller.RefreshTokens(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}