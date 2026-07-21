package categorygrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	categorygrpc "github.com/egotk/golang-advert-app/internal/features/category/controller/grpc"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_Delete(t *testing.T) {
	type deleteMockBehavior func(muc *MockuseCase, id int64)

	testTable := []struct {
		name         string
		inputRequest *categorypb.DeleteRequest
		mockBehavior deleteMockBehavior
		expectedErr  error
	}{
		{
			name:         "OK",
			inputRequest: &categorypb.DeleteRequest{Id: 1},
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Delete(gomock.Any(), id).Return(nil)
			},
		},
		{
			name:         "usecase error",
			inputRequest: &categorypb.DeleteRequest{Id: 1},
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Delete(gomock.Any(), id).Return(coreerrors.ErrNotFound)
			},
			expectedErr: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputRequest.Id)

			controller := categorygrpc.New(muc)

			_, err := controller.Delete(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
		})
	}
}
