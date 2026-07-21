package usergrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usergrpc "github.com/egotk/golang-advert-app/internal/features/user/controller/grpc"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_List(t *testing.T) {
	type listMockBehavior func(muc *MockuseCase)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	mockErr := errors.New("mock")

	testTable := []struct {
		name             string
		inputRequest     *userpb.ListRequest
		mockBehavior     listMockBehavior
		expectedResponse *userpb.UsersResponse
		expectedErr      error
	}{
		{
			name:         "OK",
			inputRequest: &userpb.ListRequest{},
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]userentity.User{
						{
							ID:          1,
							Version:     1,
							Email:       "test@example.com",
							FullName:    "Test User",
							PhoneNumber: "1234567890",
							Role:        "user",
							CreatedAt:   fixTime,
							UpdatedAt:   fixTime,
						},
					}, nil)
			},
			expectedResponse: &userpb.UsersResponse{
				Users: []*userpb.UserResponse{
					{
						Id:          1,
						Version:     1,
						Email:       "test@example.com",
						FullName:    "Test User",
						PhoneNumber: "1234567890",
						Role:        "user",
						CreatedAt:   timestamppb.New(fixTime),
						UpdatedAt:   timestamppb.New(fixTime),
					},
				},
			},
		},
		{
			name:         "usecase error",
			inputRequest: &userpb.ListRequest{},
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mockErr)
			},
			expectedErr: mockErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc)

			controller := usergrpc.New(muc)

			response, err := controller.List(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}