package favusecase_test

import (
	"context"
	"errors"
	"testing"

	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_ListIDs(t *testing.T) {
	type listIDsMockBehavior func(mr *Mockrepo, userID int64)

	mockErr := errors.New("boom")

	testTable := []struct {
		name          string
		inputUserID   int64
		mockBehavior  listIDsMockBehavior
		expectedErrIs error
		expectedIDs   []int64
	}{
		{
			name:        "OK",
			inputUserID: 1,
			mockBehavior: func(mr *Mockrepo, userID int64) {
				mr.EXPECT().ListIDs(gomock.Any(), userID).Return([]int64{1, 2, 3}, nil)
			},
			expectedIDs: []int64{1, 2, 3},
		},
		{
			name:        "repo error",
			inputUserID: 1,
			mockBehavior: func(mr *Mockrepo, userID int64) {
				mr.EXPECT().ListIDs(gomock.Any(), userID).Return(nil, mockErr)
			},
			expectedErrIs: mockErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr, testCase.inputUserID)

			uc := favusecase.New(mr)

			ids, err := uc.ListIDs(context.Background(), testCase.inputUserID)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedIDs, ids)
		})
	}
}
