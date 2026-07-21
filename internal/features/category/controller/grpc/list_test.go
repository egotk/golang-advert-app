package categorygrpc_test

import (
	"context"
	"errors"
	"testing"

	categorygrpc "github.com/egotk/golang-advert-app/internal/features/category/controller/grpc"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestController_List(t *testing.T) {
	type listMockBehavior func(muc *MockuseCase)

	parentID := int64(1)
	mockErr := errors.New("boom")

	testTable := []struct {
		name             string
		mockBehavior     listMockBehavior
		expectedResponse *categorypb.CategoriesResponse
		expectedErr      error
	}{
		{
			name: "OK",
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().List(gomock.Any()).Return(
					[]categoryentity.Category{
						{
							ID:       1,
							ParentID: nil,
							Name:     "Electronics",
						},
						{
							ID:       2,
							ParentID: &parentID,
							Name:     "Phones",
						},
					}, nil)
			},
			expectedResponse: &categorypb.CategoriesResponse{
				Categories: []*categorypb.CategoryResponse{
					{
						Id:   1,
						Name: "Electronics",
					},
					{
						Id:       2,
						ParentId: &parentID,
						Name:     "Phones",
					},
				},
			},
		},
		{
			name: "usecase error",
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().List(gomock.Any()).Return(nil, mockErr)
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

			controller := categorygrpc.New(muc)

			response, err := controller.List(context.Background(), &emptypb.Empty{})

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}
