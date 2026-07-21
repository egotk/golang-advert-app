package categoryusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/egotk/golang-advert-app/internal/core/nullable"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Patch(t *testing.T) {
	type patchMockBehavior func(mr *Mockrepo, dto categoryusecase.PatchDTO)

	newParentID := int64(2)
	invalidParentID := int64(0)
	selfParentID := int64(1)
	newName := "Updated"

	testTable := []struct {
		name          string
		inputDTO      categoryusecase.PatchDTO
		mockBehavior  patchMockBehavior
		expectedErrIs error
		expectedName  string
	}{
		{
			name: "OK set parent",
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: &newParentID},
				Name:     &newName,
			},
			mockBehavior: func(mr *Mockrepo, dto categoryusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(categoryentity.Category{
					ID:   1,
					Name: "Electronics",
				}, nil)
				mr.EXPECT().Patch(gomock.Any(), categoryentity.Category{
					ID:       1,
					ParentID: dto.ParentID.Value,
					Name:     newName,
				}).Return(nil)
			},
			expectedName: "Updated",
		},
		{
			name: "OK clear parent",
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: nil},
			},
			mockBehavior: func(mr *Mockrepo, dto categoryusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(categoryentity.Category{
					ID:       1,
					ParentID: &newParentID,
					Name:     "Electronics",
				}, nil)
				mr.EXPECT().Patch(gomock.Any(), categoryentity.Category{
					ID:       1,
					ParentID: nil,
					Name:     "Electronics",
				}).Return(nil)
			},
			expectedName: "Electronics",
		},
		{
			name: "invalid dto",
			inputDTO: categoryusecase.PatchDTO{
				ID:   0,
				Name: &newName,
			},
			mockBehavior:  func(mr *Mockrepo, dto categoryusecase.PatchDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "parent id not positive",
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: &invalidParentID},
			},
			mockBehavior:  func(mr *Mockrepo, dto categoryusecase.PatchDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "empty patch request",
			inputDTO: categoryusecase.PatchDTO{
				ID: 1,
			},
			mockBehavior:  func(mr *Mockrepo, dto categoryusecase.PatchDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "category cant be its own parent",
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: &selfParentID},
			},
			mockBehavior: func(mr *Mockrepo, dto categoryusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(categoryentity.Category{
					ID:   1,
					Name: "Electronics",
				}, nil)
			},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo GetByID error",
			inputDTO: categoryusecase.PatchDTO{
				ID:   1,
				Name: &newName,
			},
			mockBehavior: func(mr *Mockrepo, dto categoryusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(categoryentity.Category{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
		{
			name: "repo Patch error",
			inputDTO: categoryusecase.PatchDTO{
				ID:   1,
				Name: &newName,
			},
			mockBehavior: func(mr *Mockrepo, dto categoryusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(categoryentity.Category{
					ID:   1,
					Name: "Electronics",
				}, nil)
				mr.EXPECT().Patch(gomock.Any(), categoryentity.Category{
					ID:   1,
					Name: newName,
				}).Return(coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr, testCase.inputDTO)

			uc := categoryusecase.New(mr)

			category, err := uc.Patch(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedName, category.Name)
		})
	}
}
