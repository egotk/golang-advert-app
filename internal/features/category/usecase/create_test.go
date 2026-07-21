package categoryusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Create(t *testing.T) {
	type createMockBehavior func(mr *Mockrepo)

	parentID := int64(1)

	testTable := []struct {
		name          string
		inputDTO      categoryusecase.CreateDTO
		mockBehavior  createMockBehavior
		expectedErrIs error
		expectedName  string
	}{
		{
			name: "OK",
			inputDTO: categoryusecase.CreateDTO{
				ParentID: &parentID,
				Name:     "Electronics",
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, category *categoryentity.Category) error {
						category.ID = 1

						return nil
					})
			},
			expectedName: "Electronics",
		},
		{
			name: "invalid dto",
			inputDTO: categoryusecase.CreateDTO{
				Name: "",
			},
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "whitespace name",
			inputDTO: categoryusecase.CreateDTO{
				Name: "   ",
			},
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: categoryusecase.CreateDTO{
				Name: "Electronics",
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(coreerrors.ErrConflict)
			},
			expectedErrIs: coreerrors.ErrConflict,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr)

			uc := categoryusecase.New(mr)

			category, err := uc.Create(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, int64(1), category.ID)
			assert.Equal(t, testCase.expectedName, category.Name)
		})
	}
}
