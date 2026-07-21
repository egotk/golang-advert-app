package categoryusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Delete(t *testing.T) {
	type deleteMockBehavior func(mr *Mockrepo, id int64)

	testTable := []struct {
		name          string
		inputID       int64
		mockBehavior  deleteMockBehavior
		expectedErrIs error
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehavior: func(mr *Mockrepo, id int64) {
				mr.EXPECT().DeleteByID(gomock.Any(), id).Return(nil)
			},
		},
		{
			name:          "invalid id",
			inputID:       0,
			mockBehavior:  func(mr *Mockrepo, id int64) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name:    "repo error",
			inputID: 1,
			mockBehavior: func(mr *Mockrepo, id int64) {
				mr.EXPECT().DeleteByID(gomock.Any(), id).Return(coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr, testCase.inputID)

			uc := categoryusecase.New(mr)

			err := uc.Delete(context.Background(), testCase.inputID)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}
