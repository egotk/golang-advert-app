package userusecase_test

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_RefreshTokens(t *testing.T) {
	type refreshTokensMockBehavior func(mr *Mockrepo, mjs *MockjwtService)

	oldHash := corejwt.HashToken("old-refresh-token")

	testTable := []struct {
		name          string
		inputDTO      userusecase.RefreshTokensDTO
		mockBehavior  refreshTokensMockBehavior
		expectedErrIs error
		checkResult   func(t *testing.T, tokens userusecase.TokensDTO)
	}{
		{
			name:     "OK",
			inputDTO: userusecase.RefreshTokensDTO{RefreshToken: "old-refresh-token"},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetRefreshTokenByHash(gomock.Any(), oldHash).Return(userentity.RefreshToken{
					Hash:      oldHash,
					UserID:    1,
					ExpiresAt: time.Now().Add(time.Hour),
				}, nil)
				mr.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(userentity.User{ID: 1, Role: "user"}, nil)
				mjs.EXPECT().IssuePair("user", int64(1)).Return(corejwt.Pair{
					AccessToken: "new-access-token",
					RefreshToken: corejwt.RefreshToken{
						Token:     "new-refresh-token",
						IssuedAt:  time.Now(),
						ExpiresAt: time.Now().Add(720 * time.Hour),
					},
				}, nil)
				mr.EXPECT().ReissueRefreshToken(gomock.Any(), int64(1), oldHash, gomock.Any()).Return(nil)
			},
			checkResult: func(t *testing.T, tokens userusecase.TokensDTO) {
				assert.Equal(t, "new-access-token", tokens.Access)
				assert.Equal(t, "new-refresh-token", tokens.Refresh)
			},
		},
		{
			name:          "invalid dto",
			inputDTO:      userusecase.RefreshTokensDTO{RefreshToken: ""},
			mockBehavior:  func(mr *Mockrepo, mjs *MockjwtService) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name:     "token not found",
			inputDTO: userusecase.RefreshTokensDTO{RefreshToken: "unknown-refresh-token"},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetRefreshTokenByHash(gomock.Any(), corejwt.HashToken("unknown-refresh-token")).Return(userentity.RefreshToken{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
		{
			name:     "token expired",
			inputDTO: userusecase.RefreshTokensDTO{RefreshToken: "old-refresh-token"},
			mockBehavior: func(mr *Mockrepo, mjs *MockjwtService) {
				mr.EXPECT().GetRefreshTokenByHash(gomock.Any(), oldHash).Return(userentity.RefreshToken{
					Hash:      oldHash,
					UserID:    1,
					ExpiresAt: time.Now().Add(-time.Hour),
				}, nil)
			},
			expectedErrIs: userentity.ErrRefreshTokenExpired,
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

			tokens, err := uc.RefreshTokens(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkResult(t, tokens)
		})
	}
}