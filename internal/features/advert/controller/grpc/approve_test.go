package advertgrpc_test

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertgrpc "github.com/egotk/golang-advert-app/internal/features/advert/controller/grpc"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_Approve(t *testing.T) {
	type approveMockBehavior func(muc *MockuseCase, id int64)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name             string
		inputRequest     *advertpb.ApproveRequest
		mockBehavior     approveMockBehavior
		expectedResponse *advertpb.AdvertResponse
		expectedErrIs    error
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.ApproveRequest{Id: 1},
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Approve(gomock.Any(), id).Return(advertentity.Advert{
					ID:        1,
					Version:   2,
					UserID:    1,
					Title:     "Title",
					Status:    advertentity.StatusActive,
					CreatedAt: fixTime,
					UpdatedAt: fixTime,
				}, nil)
			},
			expectedResponse: &advertpb.AdvertResponse{
				Id:        1,
				Version:   2,
				UserId:    1,
				Title:     "Title",
				Status:    "active",
				CreatedAt: timestamppb.New(fixTime),
				UpdatedAt: timestamppb.New(fixTime),
			},
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.ApproveRequest{Id: 1},
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Approve(gomock.Any(), id).Return(advertentity.Advert{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputRequest.Id)

			controller := advertgrpc.New(muc)

			response, err := controller.Approve(context.Background(), testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}