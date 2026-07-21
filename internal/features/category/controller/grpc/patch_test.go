package categorygrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/nullable"
	categorygrpc "github.com/egotk/golang-advert-app/internal/features/category/controller/grpc"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	categorypb "github.com/egotk/golang-advert-app/internal/gen/category"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestController_Patch(t *testing.T) {
	type patchMockBehavior func(muc *MockuseCase, dto categoryusecase.PatchDTO)

	newParentID := int64(2)
	newName := "Updated"

	testTable := []struct {
		name             string
		inputRequest     *categorypb.PatchRequest
		inputDTO         categoryusecase.PatchDTO
		mockBehavior     patchMockBehavior
		expectedResponse *categorypb.CategoryResponse
		expectedErr      error
	}{
		{
			name: "OK set parent",
			inputRequest: &categorypb.PatchRequest{
				Id:          1,
				Name:        &newName,
				ParentIdOpt: &categorypb.PatchRequest_ParentId{ParentId: newParentID},
			},
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: &newParentID},
				Name:     &newName,
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(categoryentity.Category{
					ID:       1,
					ParentID: dto.ParentID.Value,
					Name:     *dto.Name,
				}, nil)
			},
			expectedResponse: &categorypb.CategoryResponse{
				Id:       1,
				ParentId: &newParentID,
				Name:     "Updated",
			},
		},
		{
			name: "OK clear parent",
			inputRequest: &categorypb.PatchRequest{
				Id:          1,
				ParentIdOpt: &categorypb.PatchRequest_ShouldClear{ShouldClear: &emptypb.Empty{}},
			},
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: nil},
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(categoryentity.Category{
					ID:       1,
					ParentID: nil,
					Name:     "Electronics",
				}, nil)
			},
			expectedResponse: &categorypb.CategoryResponse{
				Id:   1,
				Name: "Electronics",
			},
		},
		{
			name: "usecase error",
			inputRequest: &categorypb.PatchRequest{
				Id:   1,
				Name: &newName,
			},
			inputDTO: categoryusecase.PatchDTO{
				ID:   1,
				Name: &newName,
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(categoryentity.Category{}, coreerrors.ErrNotFound)
			},
			expectedErr: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := categorygrpc.New(muc)

			response, err := controller.Patch(context.Background(), testCase.inputRequest)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}
