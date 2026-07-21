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

func TestController_Login(t *testing.T) {
	type loginMockBehavior func(muc *MockuseCase, dto userusecase.LoginDTO)

	testTable := []struct {
		name             string
		inputRequest     *userpb.LoginRequest
		inputDTO         userusecase.LoginDTO
		mockBehavior     loginMockBehavior
		expectedResponse *userpb.LoginResponse
		expectedErr      error
	}{
		{
			name: "OK",
			inputRequest: &userpb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.LoginDTO) {
				muc.EXPECT().Login(gomock.Any(), dto).Return(
					userusecase.LoginResultDTO{
						UserID: 1,
						Tokens: userusecase.TokensDTO{
							Access:  "access-token",
							Refresh: "refresh-token",
						},
					}, nil)
			},
			expectedResponse: &userpb.LoginResponse{
				UserId:       1,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
		},
		{
			name: "usecase error",
			inputRequest: &userpb.LoginRequest{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.LoginDTO) {
				muc.EXPECT().Login(gomock.Any(), dto).Return(userusecase.LoginResultDTO{}, coreerrors.ErrUnauthorized)
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

			response, err := controller.Login(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}