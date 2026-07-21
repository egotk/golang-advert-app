package userusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_GetUserByID(t *testing.T) {
	type getUserByIDMockBehavior func(mr *Mockrepo, id int64)

	testTable := []struct {
		name          string
		inputID       int64
		mockBehavior  getUserByIDMockBehavior
		expectedErrIs error
		expectedUser  userentity.User
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehavior: func(mr *Mockrepo, id int64) {
				mr.EXPECT().GetUserByID(gomock.Any(), id).Return(userentity.User{
					ID:    1,
					Email: "test@example.com",
					Role:  "user",
				}, nil)
			},
			expectedUser: userentity.User{
				ID:    1,
				Email: "test@example.com",
				Role:  "user",
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
				mr.EXPECT().GetUserByID(gomock.Any(), id).Return(userentity.User{}, coreerrors.ErrNotFound)
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

			uc := userusecase.New(mr, NewMockjwtService(c))

			user, err := uc.GetUserByID(context.Background(), testCase.inputID)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedUser, user)
		})
	}
}