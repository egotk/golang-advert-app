package userusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_Login(t *testing.T) {
	type loginMockBehavior func(mr *Mockrepo, mjs *MockjwtService)

	correctPasswordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	futureLock := time.Now().Add(time.Hour)
	lockAfterFail := time.Now().Add(15 * time.Minute)

	testTable := []struct {
		name          string
		inputDTO      userusecase.LoginDTO
		mockBehavior  loginMockBehavior
		expectedErrIs error
		checkResult   func(t *testing.T, result userusecase.LoginResultDTO)
	}{
		{
			name: "OK",
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "correct-password",
			},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(userentity.User{
					ID:           1,
					Role:         "user",
					PasswordHash: string(correctPasswordHash),
				}, nil)
				mr.EXPECT().ResetFailedLoginCount(gomock.Any(), int64(1), int64(0)).Return(nil)
				mjs.EXPECT().IssuePair("user", int64(1)).Return(corejwt.Pair{
					AccessToken: "access-token",
					RefreshToken: corejwt.RefreshToken{
						Token:     "refresh-token",
						IssuedAt:  time.Now(),
						ExpiresAt: time.Now().Add(720 * time.Hour),
					},
				}, nil)
				mr.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResult: func(t *testing.T, result userusecase.LoginResultDTO) {
				assert.Equal(t, int64(1), result.UserID)
				assert.Equal(t, "access-token", result.Tokens.Access)
				assert.Equal(t, "refresh-token", result.Tokens.Refresh)
			},
		},
		{
			name: "invalid dto",
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "short",
			},
			mockBehavior:  func(mr *Mockrepo, mjs *MockjwtService) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "user not found",
			inputDTO: userusecase.LoginDTO{
				Email:    "unknown@example.com",
				Password: "correct-password",
			},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetUserByEmail(gomock.Any(), "unknown@example.com").Return(userentity.User{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
		{
			name: "user is locked",
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "correct-password",
			},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(userentity.User{
					ID:          1,
					Role:        "user",
					LockedUntil: &futureLock,
				}, nil)
			},
			expectedErrIs: userentity.ErrUserIsLocked,
		},
		{
			name: "wrong password",
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(userentity.User{
					ID:           1,
					Version:      2,
					Role:         "user",
					PasswordHash: string(correctPasswordHash),
				}, nil)
				mr.EXPECT().IncrementFailedLoginCount(gomock.Any(), int64(1), int64(2)).Return(nil, nil)
			},
			expectedErrIs: userentity.ErrInvalidPassword,
		},
		{
			name: "wrong password locks the user",
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(userentity.User{
					ID:           1,
					Version:      2,
					Role:         "user",
					PasswordHash: string(correctPasswordHash),
				}, nil)
				mr.EXPECT().IncrementFailedLoginCount(gomock.Any(), int64(1), int64(2)).Return(&lockAfterFail, nil)
			},
			expectedErrIs: userentity.ErrUserIsLocked,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			mjs := NewMockjwtService(c)
			testCase.mockBehavior(mr, mjs)

			uc := userusecase.New(mr, mjs)

			result, err := uc.Login(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.True(t, errors.Is(err, testCase.expectedErrIs))

				return
			}

			assert.NoError(t, err)
			testCase.checkResult(t, result)
		})
	}
}