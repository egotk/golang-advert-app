package userusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_CreateUser(t *testing.T) {
	type createUserMockBehavior func(mr *Mockrepo)

	testTable := []struct {
		name          string
		inputDTO      userusecase.CreateDTO
		mockBehavior  createUserMockBehavior
		expectedErrIs error
		checkResult   func(t *testing.T, user userentity.User)
	}{
		{
			name: "OK",
			inputDTO: userusecase.CreateDTO{
				Email:       "Test@Example.com",
				FullName:    "Test User",
				PhoneNumber: "+15551234567",
				Password:    "password123",
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, user *userentity.User) error {
						user.ID = 1

						return nil
					})
			},
			checkResult: func(t *testing.T, user userentity.User) {
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, int64(1), user.Version)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "Test User", user.FullName)
				assert.Equal(t, "+15551234567", user.PhoneNumber)
				assert.Equal(t, "user", user.Role)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")))
			},
		},
		{
			name: "invalid dto",
			inputDTO: userusecase.CreateDTO{
				Email:       "",
				FullName:    "Test User",
				PhoneNumber: "+15551234567",
				Password:    "password123",
			},
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: userusecase.CreateDTO{
				Email:       "test@example.com",
				FullName:    "Test User",
				PhoneNumber: "+15551234567",
				Password:    "password123",
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(coreerrors.ErrConflict)
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

			uc := userusecase.New(mr, NewMockjwtService(c))

			user, err := uc.CreateUser(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkResult(t, user)
		})
	}
}