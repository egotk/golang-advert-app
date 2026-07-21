package favusecase_test

import (
	"context"
	"errors"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Remove(t *testing.T) {
	type removeMockBehavior func(mr *Mockrepo)

	mockErr := errors.New("boom")

	testTable := []struct {
		name          string
		inputDTO      favusecase.RemoveDTO
		mockBehavior  removeMockBehavior
		expectedErrIs error
	}{
		{
			name: "OK",
			inputDTO: favusecase.RemoveDTO{
				AdvertID: 1,
				UserID:   1,
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().Remove(gomock.Any(), int64(1), int64(1)).Return(nil)
			},
		},
		{
			name: "invalid dto",
			inputDTO: favusecase.RemoveDTO{
				AdvertID: 0,
				UserID:   1,
			},
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo not found",
			inputDTO: favusecase.RemoveDTO{
				AdvertID: 1,
				UserID:   1,
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().Remove(gomock.Any(), int64(1), int64(1)).Return(coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
		{
			name: "repo error",
			inputDTO: favusecase.RemoveDTO{
				AdvertID: 1,
				UserID:   1,
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().Remove(gomock.Any(), int64(1), int64(1)).Return(mockErr)
			},
			expectedErrIs: mockErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr)

			uc := favusecase.New(mr)

			err := uc.Remove(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}
