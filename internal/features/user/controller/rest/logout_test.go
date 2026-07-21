package userhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/rest"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Logout(t *testing.T) {
	type logoutMockBehavior func(muc *MockuseCase)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name               string
		inputBody          string
		withClaims         bool
		mockBehavior       logoutMockBehavior
		expectedStatusCode int
	}{
		{
			name:       "OK",
			inputBody:  `{"refresh_token":"valid-refresh-token"}`,
			withClaims: true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().Logout(gomock.Any(), userusecase.LogoutDTO{
					UserID:       1,
					RefreshToken: "valid-refresh-token",
				}).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "invalid body",
			inputBody:          `invalid json`,
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			inputBody:          `{"refresh_token":"valid-refresh-token"}`,
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "usecase error",
			inputBody:  `{"refresh_token":"unknown-refresh-token"}`,
			withClaims: true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().Logout(gomock.Any(), userusecase.LogoutDTO{
					UserID:       1,
					RefreshToken: "unknown-refresh-token",
				}).Return(coreerrors.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc)

			controller := userhttp.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(testCase.inputBody))
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})

			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Logout(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)
		})
	}
}