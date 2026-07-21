package userusecase_test

import (
	"context"
	"errors"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Logout(t *testing.T) {
	type logoutMockBehavior func(mr *Mockrepo)

	mockErr := errors.New("boom")

	testTable := []struct {
		name          string
		inputDTO      userusecase.LogoutDTO
		mockBehavior  logoutMockBehavior
		expectedErrIs error
	}{
		{
			name: "OK",
			inputDTO: userusecase.LogoutDTO{
				UserID:       1,
				RefreshToken: "valid-refresh-token",
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().DeleteRefreshTokenByHash(gomock.Any(), int64(1), corejwt.HashToken("valid-refresh-token")).Return(nil)
			},
		},
		{
			name: "invalid dto",
			inputDTO: userusecase.LogoutDTO{
				UserID:       0,
				RefreshToken: "valid-refresh-token",
			},
			mockBehavior:  func(mr *Mockrepo) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: userusecase.LogoutDTO{
				UserID:       1,
				RefreshToken: "valid-refresh-token",
			},
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().DeleteRefreshTokenByHash(gomock.Any(), int64(1), corejwt.HashToken("valid-refresh-token")).Return(mockErr)
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

			err := uc.Logout(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}