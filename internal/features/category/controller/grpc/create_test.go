package categorygrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	categorygrpc "github.com/egotk/golang-advert-app/internal/features/category/controller/grpc"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestController_Create(t *testing.T) {
	type createMockBehavior func(muc *MockuseCase, dto categoryusecase.CreateDTO)

	parentID := int64(1)

	testTable := []struct {
		name             string
		inputRequest     *categorypb.CreateRequest
		inputDTO         categoryusecase.CreateDTO
		mockBehavior     createMockBehavior
		expectedResponse *categorypb.CategoryResponse
		expectedErr      error
	}{
		{
			name: "OK",
			inputRequest: &categorypb.CreateRequest{
				ParentId: &parentID,
				Name:     "Electronics",
			},
			inputDTO: categoryusecase.CreateDTO{
				ParentID: &parentID,
				Name:     "Electronics",
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.CreateDTO) {
				muc.EXPECT().Create(gomock.Any(), dto).Return(categoryentity.Category{
					ID:       1,
					ParentID: dto.ParentID,
					Name:     dto.Name,
				}, nil)
			},
			expectedResponse: &categorypb.CategoryResponse{
				Id:       1,
				ParentId: &parentID,
				Name:     "Electronics",
			},
		},
		{
			name: "usecase error",
			inputRequest: &categorypb.CreateRequest{
				Name: "Electronics",
			},
			inputDTO: categoryusecase.CreateDTO{
				Name: "Electronics",
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.CreateDTO) {
				muc.EXPECT().Create(gomock.Any(), dto).Return(categoryentity.Category{}, coreerrors.ErrConflict)
			},
			expectedErr: coreerrors.ErrConflict,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := categorygrpc.New(muc)

			response, err := controller.Create(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}
