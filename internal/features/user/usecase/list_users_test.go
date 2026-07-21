package userusecase_test

import (
	"context"
	"errors"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_ListUsers(t *testing.T) {
	type listUsersMockBehavior func(mr *Mockrepo)

	negativeLimit := int64(-1)
	negativeOffset := int64(-1)
	mockErr := errors.New("boom")

	testTable := []struct {
		name          string
		inputLimit    *int64
		inputOffset   *int64
		mockBehavior  listUsersMockBehavior
		expectedErrIs error
		expectedUsers []userentity.User
	}{
		{
			name:        "OK",
			inputLimit:  nil,
			inputOffset: nil,
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().ListUsers(gomock.Any(), (*int64)(nil), (*int64)(nil)).Return(
					[]userentity.User{{ID: 1, Email: "test@example.com", Role: "user"}}, nil)
			},
			expectedUsers: []userentity.User{{ID: 1, Email: "test@example.com", Role: "user"}},
		},
		{
			name:          "negative limit",
			inputLimit:    &negativeLimit,
			inputOffset:   nil,
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name:          "negative offset",
			inputLimit:    nil,
			inputOffset:   &negativeOffset,
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name:        "repo error",
			inputLimit:  nil,
			inputOffset: nil,
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().ListUsers(gomock.Any(), (*int64)(nil), (*int64)(nil)).Return(nil, mockErr)
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

			uc := userusecase.New(mr, NewMockjwtService(c))

			users, err := uc.ListUsers(context.Background(), testCase.inputLimit, testCase.inputOffset)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedUsers, users)
		})
	}
}