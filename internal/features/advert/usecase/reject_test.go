package advertusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Reject(t *testing.T) {
	type rejectMockBehavior func(mr *Mockrepo, id int64)

	testTable := []struct {
		name          string
		inputID       int64
		mockBehavior  rejectMockBehavior
		expectedErrIs error
		expectedID    int64
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehavior: func(mr *Mockrepo, id int64) {
				mr.EXPECT().ChangeStatus(gomock.Any(), id, advertentity.StatusInitial, advertentity.StatusRejected).Return(
					advertentity.Advert{ID: 1, Status: advertentity.StatusRejected}, nil)
			},
			expectedID: 1,
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
				mr.EXPECT().ChangeStatus(gomock.Any(), id, advertentity.StatusInitial, advertentity.StatusRejected).Return(
					advertentity.Advert{}, coreerrors.ErrNotFound)
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

			uc := advertusecase.New(NewMockfavRepo(c), mr, NewMockstorage(c))

			advert, err := uc.Reject(context.Background(), testCase.inputID)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedID, advert.ID)
		})
	}
}
