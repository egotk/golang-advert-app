package usergrpc_test

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	usergrpc "github.com/egotk/golang-advert-app/internal/features/user/controller/grpc"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_GetByID(t *testing.T) {
	type getByIDMockBehavior func(muc *MockuseCase, id int64)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name             string
		inputRequest     *userpb.GetByIDRequest
		mockBehavior     getByIDMockBehavior
		expectedResponse *userpb.UserResponse
		expectedErr      error
	}{
		{
			name:         "OK",
			inputRequest: &userpb.GetByIDRequest{UserId: 1},
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().GetUserByID(gomock.Any(), id).Return(userentity.User{
					ID:          1,
					Version:     1,
					Email:       "test@example.com",
					FullName:    "Test User",
					PhoneNumber: "1234567890",
					Role:        "user",
					CreatedAt:   fixTime,
					UpdatedAt:   fixTime,
				}, nil)
			},
			expectedResponse: &userpb.UserResponse{
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
		{
			name:         "usecase error",
			inputRequest: &userpb.GetByIDRequest{UserId: 1},
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().GetUserByID(gomock.Any(), id).Return(userentity.User{}, coreerrors.ErrNotFound)
			},
			expectedErr: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputRequest.UserId)

			controller := usergrpc.New(muc)

			response, err := controller.GetByID(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}