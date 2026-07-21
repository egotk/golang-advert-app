package userhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/rest"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_RefreshTokens(t *testing.T) {
	type refreshTokensMockBehavior func(muc *MockuseCase, dto userusecase.RefreshTokensDTO)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDTO             userusecase.RefreshTokensDTO
		mockBehavior         refreshTokensMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"refresh_token":"old-refresh-token"}`,
			inputDTO: userusecase.RefreshTokensDTO{
				RefreshToken: "old-refresh-token",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.RefreshTokensDTO) {
				muc.EXPECT().RefreshTokens(gomock.Any(), dto).Return(
					userusecase.TokensDTO{
						Access:  "new-access-token",
						Refresh: "new-refresh-token",
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`,
		},
		{
			name:               "invalid body",
			inputBody:          `invalid json`,
			mockBehavior:       func(muc *MockuseCase, dto userusecase.RefreshTokensDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "usecase error",
			inputBody: `{"refresh_token":"expired-refresh-token"}`,
			inputDTO: userusecase.RefreshTokensDTO{
				RefreshToken: "expired-refresh-token",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.RefreshTokensDTO) {
				muc.EXPECT().RefreshTokens(gomock.Any(), dto).Return(userusecase.TokensDTO{}, coreerrors.ErrUnauthorized)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := userhttp.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(testCase.inputBody))
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.RefreshTokens(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}